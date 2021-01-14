from mininet.topo import Topo
from mininet.net import Mininet
from mininet.node import Controller
from mininet.cli import CLI
from mininet.link import TCLink


def multipathTopo():
    net = Mininet(build = False,link=TCLink)
    
    c0 = net.addController("controller")

    server = net.addHost('server')
    client = net.addHost('client')
    router = net.addHost('router')

    leftSwitch = net.addSwitch('s1')
    rightSwitch = net.addSwitch('s2')
    serverSwitch = net.addSwitch('s3')

    net.addLink(client, leftSwitch)  		
    net.addLink(client, rightSwitch)

    net.addLink(router, leftSwitch)		
    net.addLink(router, rightSwitch)
    net.addLink(router, serverSwitch) 
    net.addLink(serverSwitch, server)

    net.build()
	
    server.cmd('./mininet-topology/server-config.sh')
    client.cmd('./mininet-topology/client-config.sh')
    router.cmd('./mininet-topology/router-config.sh')

    controller = net.controllers[0]
    controller.start()
    
    leftSwitch.start([controller])
    rightSwitch.start([controller])
    serverSwitch.start([controller])

    return net

if __name__ == "__main__":
    net = multipathTopo()
    CLI(net)

