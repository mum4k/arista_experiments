#!/usr/bin/python
"""Fetches network interfaces from an Arista box."""

import pyeapi

def main():
  pyeapi.load_config('eapi.conf')
  node = pyeapi.connect_to('arista')

  ifaces = node.api('interfaces')
  print(ifaces.getall())

if __name__ == '__main__':
    main()
