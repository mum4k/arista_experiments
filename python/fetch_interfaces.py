#!/usr/bin/python
"""Fetches network interfaces from an Arista box.

If the IF_NAME interface doesn't have its primary IP set to IP_ADDRESS,
it will be configured.
"""

import pyeapi

IF_NAME='Management1'
IP_ADDRESS='10.1.2.90/24'


class Error(Exception):
  """Error class for this module."""


class NotFoundError(Error):
  """Returned when the specified interface cannot be found."""


def Connect():
  """Connects to the Arista router.
  
  Returns: pyeapi.Node, the connected router.
  """
  pyeapi.load_config('eapi.conf')
  return pyeapi.connect_to('arista')

def ShowInterfaces(node):
  """Executes the 'show interfaces' command.

  Args:
    node: pyeapi.Node, the target router.

  Returns:
    All interfaces in a dict object.
  """
  ifaces = node.api('interfaces')
  return ifaces.getall()

def HasPrimaryIp(node, name, ip):
  """Determines if the interface has the provided primary IP.

  Args:
    node: pyeapi.Node, the target router.
    name: string, name of the interface.
    ip: string, the IP address to check for.

  Returns:
    A bool, True if the interface has the IP address.

  Raises:
    NotFoundError: When the router doesn't have the specified interface.
  """
  iface = node.api('ipinterfaces').get(name)
  if iface is None:
    raise NotFoundError('interface %s not found' % (name,))
  return iface['address'] == ip


def SetPrimaryIp(node, name, ip):
  """Sets the primary IP address on the interface.

  Args:
    node: pyeapi.Node, the target router.
    name: string, name of the interface.
    ip: string, the IP address to set.

  Returns:
    A bool, True if the operation succeeded.

  Raises:
    NotFoundError: When the router doesn't have the specified interface.
  """
  return node.api('ipinterfaces').set_address(name, ip)


def main():
  node = Connect()
  print('ShowInterfaces:\n%s\n\n' % ShowInterfaces(node))

  if HasPrimaryIp(node, IF_NAME, IP_ADDRESS):
    print('Interface %s already has IP %s.\n' % (IF_NAME, IP_ADDRESS))
    return

  print('Adding IP %s to interface %s.\n' % (IP_ADDRESS, IF_NAME))
  if not SetPrimaryIp(node, IF_NAME, IP_ADDRESS):
      print('failed to configure IP %s on interface %s.\n' % (IP_ADDRESS, IF_NAME))


if __name__ == '__main__':
    main()
