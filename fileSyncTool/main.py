from file_sync.server import server
from file_sync.file_sync import initialize


if __name__ == "__main__":
    import argparse
    import sys

    parser = argparse.ArgumentParser(description='File synchronization Tool, used for chatting.')

    update_cmd_parser = parser.add_argument_group('update')
    update_cmd_parser.add_argument('--file', type=str, help='File path to sync')

    server_cmd_parser = parser.add_argument_group('server')
    server_cmd_parser.add_argument('--start', action='store_true', help='Start the chatting server')
    server_cmd_parser.add_argument('--stop', action='store_true', help='Stop the chatting server')

    if len(sys.argv) < 2:
        print("Sub command required")
        sys.exit(1)

    args = parser.parse_args(sys.argv[2:])

    # Choose which sub-command
    # if sys.argv[1] == 'update':
    #     if args.file is None:
    #         print("usage of update:")
    #         parser.print_help()
    #         sys.exit(1)
    #     update(args.file)
    # elif sys.argv[1] == 'uninstall':
    #     uninstall()
    if sys.argv[1] == 'server':
        if args.start:
            server(True)
        elif args.stop:
            server(False)
        else:
            print("usage of server:")
            parser.print_help()
            sys.exit(1)
    elif sys.argv[1] == 'init':
        initialize()
    else:
        print("Unknown sub command")
        parser.print_help()
        sys.exit(1)

