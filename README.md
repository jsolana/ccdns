# ccdns

A DNS Forwarder is a nameserver used to resolve DNS queries instead of directly using the authoritative nameserver chain. Often they are used to sit on the edge of a local area network and provide DNS resolution to the computers on the network, reducing external traffic and speeding up external access by serving the answer from a local cache.

## Step Zero

In this introductory step you’re going to set your development environment up ready to begin developing and testing your solution.

I’ll leave you to choose your target platform, setup your editor and programming language of choice. I’d encourage you to pick a tech stack that you’re comfortable doing network programming with - we’re going to be creating, sending and receiving UDP packets.

## Step 1

In this step your goal is to create a UDP server that will listen on a specified port for incoming requests. Your server should use port 53 if no port is provided, otherwise it should use the provided port (i.e. add a command line option to specify it). Note for testing we’ll use a port above 1023 as many OSs don’t allow non-privileged users to use a port below 1023.

Once you have your server running you can test it with dig. That might look like this:

```bash
% dig @127.0.0.1 -p 1053 www.google.com
``````

With the server so far running in another terminal.

```bash
% ccdns
MSG received
```

All the server does so far is listen for messages, so I’m just printing out the fact that something was received. We’ll handle it in a later step.

## Step 2

In this step your goal is to parse the request packet that has been sent to your server. A DNS message has:

- A header.
- A questions section.
- An answer section.
- An authority section.
- An additional section.

The header is defined in [RFC 1035 Section 4.1.1](https://datatracker.ietf.org/doc/html/rfc1035?utm_source=substack&utm_medium=email#section-4.1.1)

A question section has three fields: the name that we’re looking up, a record type and the class. The name is a byte string, the other two fields are 1b-bit integers.

The question is defined in the same RFC, [Section 4.1.2](https://datatracker.ietf.org/doc/html/rfc1035?utm_source=substack&utm_medium=email#section-4.1.2).

The name field is an encoded version of the host name string, so if we’re looking up [dns.google.com](https://dns.google/) it will be encoded as the following (bytes) 3dns6google3com0.

The format of the answer, authority, and additional sections is defined in [Section 4.1.3 of the RFC](https://datatracker.ietf.org/doc/html/rfc1035?utm_source=substack&utm_medium=email#section-4.1.3).

I would suggest using [test-driven development](https://en.wikipedia.org/wiki/Test-driven_development) (TDD) to build tests for code to serialise and deserialise these requests.

Once you’ve done that try your code with dig again and see if you can print out a human readable interpretation of the request. For example:

```bash
Header - id:40500 flags:288 questions:1 answers:0 authorities:0 additionals:1 
Questions - name:www.google.com type:1 class:1 
Answers: Empty
Authorities: Empty
...
```

For the query:

```bash
% dig @127.0.0.1 -p 1053 www.google.com
```

## Step 3

In this step your goal is to forward the request to a DNS server to actually resolve the request. For this challenge we’ll use [Google’s public DNS server](https://developers.google.com/speed/public-dns) as an example (8.8.8.8).

Your forwarder should receive the request, unpack it and then re-pack it for forwarding to a DNS server to be resolved. So your code should be receiving a request, and sending the same request to 8.8.8.8 port 53. For this step it’s enough to print out that you’ve forwarded it.

```bash
% ccdns
Forwarded Question
```

## Step 4

In this step your goal is to receive an answer from the DNS server, unpack it and forward the answer to the original client that asked for it.

To do this you’ll need to keep a record of which client asked which question, then when you get an answer to a question, look up the question and send the answer to the client that asked it.

Once this step is complete you should be able to send a request to your forwarder and see a response like so:

```bash
% dig @127.0.0.1 -p 1053 www.google.com

; <<>> DiG 9.10.6 <<>> @127.0.0.1 -p 1053 www.google.com
; (1 server found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 43712
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 512
;; QUESTION SECTION:
;www.google.com.                        IN      A

;; ANSWER SECTION:
www.google.com.         118     IN      A       142.250.179.228

;; Query time: 11 msec
;; SERVER: 127.0.0.1#1053(127.0.0.1)
;; WHEN: Wed Jan 10 18:09:54 GMT 2024
;; MSG SIZE  rcvd: 59
```

## Step 5

In this step your goal is to cache the results of successful DNS resolution and respond to future requests with the same question from the local cache. The cached value should only be served if it is still within the TTL provided in the answer section, otherwise a new lookup should be done.

Test that your server only forwards requests that are not cached or for which the TTL has expired and once done, congratulations you’ve built a DNS Forwarder.
