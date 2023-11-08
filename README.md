# sleeveless
simulation-based testing framework for kubernetes control plane components

## Background
### Concepts
- [The Mechanics of Kubernetes](https://dominik-tornow.medium.com/the-mechanics-of-kubernetes-ac8112eaa302)
    - good high level overview of how kubernetes "works"

- [Optimistic Concurrency](https://www.learnsteps.com/what-is-optimistic-concurrency-how-does-is-it-help-to-scale-the-kubernetes-cluster/)
- [how etcd works with Kubernetes](https://learnk8s.io/etcd-kubernetes)

### Blog Posts
- [Events: the DNA of Kubernetes](https://www.mgasch.com/2018/08/k8sevents/#fn:8)
- [The Anatomy of a Kubernetes ListWatch(): Prologue](https://www.mgasch.com/2021/01/listwatch-prologue/#fnref:4)
- [Onwards to the Core: etcd](https://www.mgasch.com/2021/01/listwatch-part-1/)
- [Kubernetes API Server Adventures: Watching and Caching](https://danielmangum.com/posts/k8s-asa-watching-and-caching/)

### Papers
- [Reasoning about modern datacenter infrastructures using partial histories (HotOS '21)](https://sigops.org/s/conferences/hotos/2021/papers/hotos21-s11-sun.pdf)
    - a short workshop paper on why Kubernetes is hard. The 3 articles in the "Blog Posts" section above describe in more detail the programming environment that this paper is talking about.

- [Automatic Reliability Testing for Cluster Management Controllers (OSDI '22)](https://www.usenix.org/system/files/osdi22-sun.pdf)
    - This is the paper for the tool "Sieve" that has inspired my current research a great amount. You can think of it as a sort of fault injection testing tool that automatically finds fault tolerance bugs in individual Kubernetes controllers.

- [What goes wrong in serverless runtimes? A survey of bugs in Knative Serving (SESAME '23)](https://dl.acm.org/doi/abs/10.1145/3592533.3592806)
    - paper that I wrote earlier this year with Andrew and Lindsey studying the Knative platform, which is a serverless platform built on top of Kubernetes. "Serverless" is a somewhat overloaded word these days, but I like to think of "serverless platforms" as platforms that are able to decouple computation from storage or networking resources such that the platform can manage these independently. By doing so, you can support "serverless functions" or services that only spin up when a request to them is made, and can go away when they aren't in use. To achieve this decoupling, your infrastructure has to have a lot of dynamic moving parts. Knative does this by implementing a variety of controller components that coordinate with eachother in complex ways. Of course, there have been bugs in the ways these controllers interact, and we talk more about them in this paper. These types of bugs are what we're trying to build a tool to detect!



