**Group Number - 10**

*[Saptarsi Saha](https://github.com/saptarsi96)*           *[Manisha Jhunjhunwala](https://github.com/ManishaJhunjhunwala)*          *[Shailesh Navale](https://github.com/Shailesh122)*         *[Madhur Navandar](https://github.com/mnavandar)*  
*2020H1030109*         *2020H1030125*                      *2020H1030169*          *2020H1030163*

**Abstract**

In the current scenario, QUIC(Quick UDP Internet Connection), developed by Google is being used extensively in Chrome to reduce page-load times, rebuffers and to have a low latency of client-server communication. It has replaced the traditional HTTPS/TLS/TCP stack and has included some cryptographic mechanisms to authenticate the server along with leveraging the reliability and congestion control mechanisms of modern TCP stacks. But it cannot expolit the multiple paths and several wireless interfaces existing between today's mobile devices. Here comes the multipath version of QUIC - MPQUIC which can take advantage of multiple paths existing between the two communicating hosts simultaneously resulting in increased utilisation of resources and efficiency. Although this protocol has not been implemented in the real-world setups, we have tried to simulate the working of MPQUIC over a mininet VM and have conducted some file transfer experiments using MPQUIC(built over UDP). Also, in the internet, web pages can be rendered as the multiple streams for e.g a web page might render the streams collecting the HTML text, video, audio, advertisement,etc. The implementation of the MPQUIC as presented in the Conext`17 research paper is not stream aware. So, we have studied the Priority Bucket scheduler([PStream](https://dl.acm.org/doi/abs/10.1145/3396851.3402923)) and have presented our ideas to improvise on the time complexity of the scheduling algorithm from O($n^2$) to O($nlogn$).

#### *I. Introduction*

In this paper we are presenting the implementation of the MPQUIC using GO language. We were able to implement the MPQUIC using the open-source QUIC-Go Repository by Quentin De Coninck.  
The first part of the project deals with the scenario of transferring/uploading a large file between the client and server in a multipath network setup. We have been successful in implementing the transfer of files of multiple sizes over different network conditions and characteristics(bandwidth, delay, losses) of the disjoint paths. We have also tested the resiliency of our multipath topology by dynamically making one of the interfaces down.  
For the second part of the project, we have studied the Priority Bucket scheduler([PStream](https://dl.acm.org/doi/abs/10.1145/3396851.3402923)) for the MPQUIC which schedules the data as streams in O($n^2$) time complexity based on the priorities of the stream. By scheduling the streams based on the priority we reduce the response time for the web page rendering. We have tried to implement the PStream scheduler by referring to the research paper's source but could only proceed partially due to several certificate errors.  
In the third part, we have analysed the existing PStream algorithm extensively and it seems that we can fill some missing gaps in it. So we have proposed a improvisation in the current algorithm and proved our idea to reduce the time complexity to O(nlogn), with supporting calculations.
***
```
		Part A: Demonstration of File transfer over MPQUIC
```
In this part we will have a quick look on how we can transfer files from host to server using MPQUIC. We have first tried to configure MPQUIC protocol in the real environment but was unable to do so because to support MPQUIC protocol we need to configure multiple tools to build and test our application as it is quite new and it’s module is not yet deployed. We are focused on transmitting single files through multiple paths and to observe the difference in time it takes when we schedule a packet through a single path(as it was done previously in gQUIC) and through multipath QUIC.

#### *II. MPQUIC Setup*

We have decided to go with the decision to use virtual machines to host all our test and then transfer the results. We could have gone with installing mininet directly in Ubuntu and proceed, but that would have needed a lot of tools to configure, build and test, because both QUIC and MPQUIC are still in development phase and are not built as modules yet.  
We have decided to use Mininet VM available in the [CoNEXT 2017 Artifacts](https://multipath-quic.org/2017/12/09/artifacts-available.html) for experimentation purposes. It comes bundled with Minitopo, a very useful helper tool to build multipath networks over Mininet.

#### *III. Simulating Networks*

We decided to simulate a multi-path network using  
the network simulator tool mininet,which was installed inside  
a Linux virtual machine.  
After installation of Mininet VM we first created the topology by using the mininet python scripts and/or minitopo sample.topo file which we designed to connect a host to a router using 2 paths.

Our created topology as following:

![topo_final.png](:/d360b06037bb4ce997a7d692aefc5886)  
*Fig 1. Two disjoint paths between Client and Server*

The network has two hosts with client having two interfaces  
and server having one, and two switches linking the two hosts. The switches are then connected to a router which connects to the server. The IP addresses  
assigned to the hosts are as follows  
**Client**:  
H1-eth0 : 10.0.0.1  
H1-eth1 : 11.0.0.1  
**Server** : 100.0.0.1

The topology creation was a little bit tricky because with our initial setup, the traffic was not evenly distributed between the 2 links. So we have made some tweaks to utilise all the paths equally.

**Minitopo to the rescue:** After careful consideration, we decided to move on to minitopo, which is a simple tool, based on mininet, to boot a simple network with n paths and run experiments between two hosts. We ran multiple experiments by tweaking the parameters in our sample.topo file as described in the later parts of the report.

#### *IV. File Transfers over MPQUIC(UDP) protocol*

After successfully transferring files through QUIC we moved our focus to MPQUIC for which we required a parameter : **CreatePaths** which was missing in [lucas-clemente’s QUIC-GO](https://github.com/lucas-clemente/quic-go) implementation. We referred to the CoNEXT 2017 paper by Quentien De Conick which comes bundled with all the required tools such as minitopo and has provision to use this parameter too.

Steps that we performed while transferring file from client to server using multiple path:

1.  Create a topology in mininet as in Fig 1.
    
2.  As quic-go implementation is only based on single path connection to transfer file packets between client and server in order to make it multi-path we needed to change the code a little bit and added the following snippet in both client and server go file.  
    `quicConfig := &quic.Config { CreatePaths: true, }`
    
3.  The file transfer code for both the sides(client and server) - server.go and client.go files is firstly built inside mininet VM.
    
4.  The server and client is run in seperate xterms.
    
5.  Server waits for the client to initiate a connection request. When the connection is established, the server decides whether to use multiple paths or not based on the file-size. We have provided a threshold in the server side, if the file size is bigger than the threshold of 10KB, multiple streams are created and transfer happens using MPQUIC.
    
6.  We have attached a storage folder for server while running the client so that the uploaded file gets stored there.
    

We observe that we were able to send the file successfully as shown in fig 2 & 3 below but utilization of both path is not efficient that we can see in fig 4.

So, how did we check that the tranfer was using MPQuic protocol only and not TCP. We checked the wireshark traces with these endpoints and confirmed that all these packets were using the UDP protocol.  
<img src=":/ac5bb17c296b4676904ec6ebef0ecf7a" alt="client.png" width="311" height="164" class="jop-noMdConv">       <img src=":/dce4a31d6f6641cfa07ceeb5b38558cd" alt="server.png" width="325" height="166" class="jop-noMdConv">

Fig 2: Client Xterm:                                                    Fig 3: Server Xterm:

<img src=":/467a54e40fe046bb8685c255e9e2d9e6" alt="uneven_transfer_stats.png" width="258" height="61" class="jop-noMdConv">    <img src=":/e5a26ad41da442d9ac940d1197df41a9" alt="even_transfer.png" width="374" height="59" class="jop-noMdConv">

*Fig 4: Uneven traffic through 2 paths*       *Fig 5: Traffic between 2 paths is evenly distributed(appx.)*

In order to efficiently use both the paths in topology, we changed the code in [config-client.sh](http://config-client.sh) file

`- route add default gw 10.0.0.2 client-eth0`  
`+ ip route add default scope global nexthop via 10.0.0.2 dev client-eth0`

We tried running the same experiment with a sample minitopo file and utilizing 4 channels for communication and got the same approx result whereby the packets were being divided equally amongst the channels and the links were being used effectively.

<img src=":/3208919a7d0240d0ba1d080ba399c5e7" alt="4 channels.png" width="408" height="53" class="jop-noMdConv">     

*Fig 6: Server IP is 10.1.0.1(in violet) and the client is connected to 4 interfaces(10.0.0.1,10.0.1.1,10.0.2.1,10.0.3.1) in red.*

**Initial bwm-ng results also confirm the same.**

<img src=":/b4b36bcec6f64e3f8844edf8faefcf58" alt="bwmng_4channels.png" width="381" height="215" class="jop-noMdConv">

#### *V. Checking the resiliency of the multipath network*

As the concept of multi-path network says, that if any path or stream goes down due to changing network connections, there should not be any hindrance with the connection if atleast one path is up.  
All the traffic should get distributed in the left-over paths which are active.

But we discovered that our topology was not able to handle it at first. We tried several ways to enable our code to make this choice dynamically but we faced many challenges in making this setting. We tried to tinker with the topology but discovered that there was some issue with the mininet which didn't allow us this functionality. So whenever we made any path down, we got a segmentation fault.  
We thought over it and decided to test in another way, that whether our code needed any change or it was really an issue with the mininet VM setup.

We configured 2 mininet VM's, one for client and another for server. Server having one interface, but for the client side, we tried to add another interface by making changes to the NAT settings.  
We then performed the whole experiment of transferring file again, this time by adding 2 other steps.

1.  We made one interface down before running the client.
2.  We made one interface down between the transfer on the fly.

We were surprised to see the results. So with the same client and server code and the topo setup, we could dynamically transfer the traffic to the left over paths if we made any interface down.

The wireshark traces confirmed the same.



<img src=":/03b332d478244b42bbd860694d826707" alt="4 channels.png" width="265" height="57" class="jop-noMdConv">       <img src=":/93eadf0923534cafb09c762722f0d9eb" alt="4 channels.png" width="392" height="55" class="jop-noMdConv">

Fig 7:                                                             Fig 8:

*Fig 7.One interface is down and the traffic is routed through the other interface. Server IP is 192.168.221.145 & client is 192.168.221.150*

*Fig 8. We cut interface with IP 192.168.221.146 dynamically while file transfer is happening and can see that the traffic gets routed through the other interface.*

#### *VI. Challenges Faced*

This project was not short of any challenges, from modifying the QUIC-go code to run for MPQUIC to changing the topology files for different scenarios and getting stuck at the expired certificates, we faced it all. Some of them are listed below:

**1\. Implementing file-transfer in QUIC.**  
While using quic-go client the method was set to use get as default. We modified it to use post method and implemented the quic-protocol.  
`- rsp, err := hclient.Get(addr)`  
`+ rsp, err := hclient.Post(addr, "binary/octet-stream", file)`

**2\. Multipath Next!**  
To implement multipath-quic we had to make changes to the existing code and include the  
`quicConfig := &quic.Config{ CreatePaths: true, }` parameter. This parameter is not accessible by the quic-go module and we had to fall-back to MPQUIC-VM.

**3\. Design Changes Finally.**  
We were using 2 paths and the traffic ratio was divided in 40:1.  
We found critical issues in our topology and made changes in the client-mininet topology to bring down the traffic ratio to 1:1 (approximately).

#### *VII. Testing with different network conditions*

We have used tools like ***Netem*** and ***bwm-ng*** to modify the bandwidth, delay and losses of all the paths and observe the results.

<img src=":/d5a142e9f34d4052aa9f079f36c9887e" alt="Number of Paths time comparison.jpeg" width="248" height="158" class="jop-noMdConv">                          <img src=":/dd583fb875ca47759eb8ebe6437e1952" alt="Bandwidth comparison_10.0.0.1 thke clients.jpeg" width="253" height="164" class="jop-noMdConv">  
Fig: Number of Paths versus Time comparison            Fig: Bandwidths versus transferred bytes

So, in Fig 1: We can observer that MPQUIC has taken a considerably less time than QUIC for large file transfers. Also, when we increased the number of paths, the transmission time kept on decreasing.

In Fig 2, we have configured the 5 interfaces of client with bandwidths in the ratio of 1:2:4:8:16. We can see that on paths with higher bandwidths, number of bytes transferred was more.

* * *

```
					PART B: Priority stream bucket

```
For this part, we have read and presented the paper on [4][PriorityBucket: A Multipath-QUIC Scheduler on Accelerating First Rendering Time in Page Loading](https://dl.acm.org/doi/abs/10.1145/3396851.3402923). We have also tried to implement it but faced some difficulties which are mentioned below.
The response time of a web page can be the time a web server takes to show the first visual feedback to the user. Visual feedback is the HTML, CSS, or Javascript components that make up the webpage. So the idea is to prioritize such streams that make up the front-end design of the web page over the other streams such as advertisement.

Protocols like HTTP 1.1 or HTTP 2.0 suffer from the problem of HOL blocking protocol since they were using TCP. To remove the unnecessary blockage of independent requests, QUIC protocol multiplexes HTTP requests/responses over a single UDP connection. With UDP protocol QUIC ensures that late packets from each individual stream will not block other concurrent streams which resolves HOL blocking problems.

In this project we will focus on an algorithm which deals with scheduling a stream to different paths based on path bandwidth or number of higher priority streams in the path and number of equal priority streams in a path compared to a stream to be scheduled by using MPQUIC protocol. The idea behind such scheduling of the packets over different paths such that the completion time of the stream over each scheduled path is equal. In this way we are trying to optimize the completion time and complete the transfer in more efficient manner.

As show in the below figure the buckets with lower number will have the higher priority. Inside the bucket scheduling will take place in round robin fashion.

![f258a803e524ca3c5843ae880836d8f2.png](:/afa7e8bda3bc401dba844b2ec4ea6dec)

Fig 1: Structure of priority bucket design [4]

The streams with higher priority can preempt the stream with lower priority during the transfer.  
![f56811f789827a2f8fc5f92fcfd9bd13.png](:/fb4f29374f4f4f7f860ef40e33004663)  
**Stream scheduling approach:**

Below table, represents the symbols and meanings of the terms used in this algorithm.

![7d1d5e36da0ba67d0947213b6cbcd744.png](:/00ac955936d64167b45a5a76c95d82b3)

Along with the parametrers mentioned in table 1, algorithm introduces one more parameter $r_i$.  
$r_i$ = ${h_i + e_i}\over{b_i}$ \+ $o_i$

$\therefore r_i$ = time required for transmission of equal and higher priority streams + one time delay.

i.e the time required for completion if we do not send any data over path i.

It can be divided in the below steps:

**Step 1:**

Sort the paths in ascending order of $r_i$.

**Step 2:**

The fillgap process schedules data to balance the completion time gap of multiple paths. It fetches two paths with smallest $r_i$ and will try to schedule balance out the time.e.g. At first the algorithm will consider path1 and path2 with smallest $r_i$ i.e. $r_1$ and $r_2$.

We want to find $d_1$ such that it satisfies below equation.

$r_1$ \+ ${d_1}\over{b_1}$ = $r_2$ (here $b_1$ is the bandwidth of path1 )

On solving,

$d_1$ = $(r_2 - r_1)$ x $b_1$

Similarly this process will be repeated until all the data is allocated or completion time of the paths becomes equivalent. This process is called as gap filling process in research paper.

**Step 3**

If any data is still not have been scheduled, schedule it in the ratio of the bandwidths of the paths.

**Implementation Attempt**

In order to implement this research paper, we have tried to use already [implemented github](https://github.com/Xiang-Shi/PStream) repo from one of the author Xiang-shi. We PStream repository in the local virtual machine and tried to set it up.  
Initially, the repo was not able to run due to the path issues. In order to implement this we had to clone this repository in the repository used implementation of Part A. We also had to change the name of this repository from 'PStream' to 'pstream'.

When we tried to run this we have faced the few errors which we were not able to resolve yet.

<img src=":/a83dc5a6e7c840c390f52a0824d7c65b" alt="210beb104458c334c87badfe358709d0.png" width="299" height="178" class="jop-noMdConv">    <img src=":/0bf6596faec04c7e96754247178b4276" alt="56171b630ca8454f7d8c0fafff887755.png" width="330" height="174" class="jop-noMdConv">

Fig 2:Pstream Client                                                     Fig 3:Pstream Server

As show in the figure 2 and 3, we faced the issue with the certificate and we will try to resolve this in future.
* * *

**References:**

1.  [https://multipath-quic.org/2017/12/09/artifacts-available.html](https://multipath-quic.org/2017/12/09/artifacts-available.html)
2.  [https://github.com/qdeconinck/mp-quic](https://github.com/qdeconinck/mp-quic)
3.  [https://github.com/lucas-clemente/quic-go](https://github.com/lucas-clemente/quic-go)
4.  [PriorityBucket: A Multipath-QUIC Scheduler on Accelerating First Rendering Time in Page Loading](https://dl.acm.org/doi/abs/10.1145/3396851.3402923)
5.  [https://github.com/Xiang-Shi/PStream](https://github.com/Xiang-Shi/PStream)
