# 2/25

Attending:
E: Ed (@ceezy)
I: Ian (@wesspacito)
C: Cam (@formomosan)

Notes taken by: Cam

### Deployments

- E: Kubernetes + Google Cloud
  - Good tool, abstracts away
  - Not too difficult but need to learn the workflow
- C: Currently set up with Docker + AWS + Compose
- E: We can use Terraform
  - Encode infrastructure without clicking around, ssh, etc.

Conclusion: E will look into Kubernetes + Google Cloud

### Naming + Domains + Where we fit in

- I: Go terms
  - gravitated towards `tengen.io`
- E:
  - Western go community:
    1. Baby boomers w/minimal tech experience
    2. younger, nerdy crowd, influenced by anime
  - Community responds well to Japanese terms
- I: Western audience disconnected with Go?
- E: Asian servers are heavily siloed and are _super_ gnarly
  - Ex: Fox
  - Many have betting, micro transactions, etc.
  - Some have Go problems but force users to pay to get answers
  - Traditionally it's been hard to break into the Asian audience
- I: Why not focus on growing Western community?
- E:
  - Ex: OGS & KGS
  - Lacking social tools
  - Fewer strong players on Western servers
  - Better than 1 or 2 dan, you won't find many games on OGS
- I: A very clean experience will bring together a dedicated, newer community
- E:
  - Existing tutorials won't lead you into games, not seamless
  - How do we lead naturally into games without being stressed out? There's a market there

### Feature development

- I: Game Clocks
- E: variants of Go are not MVP but very interesting
- C: **Board Size**
- E:
  - **Game Review tools**
  - learning resources
    - Import/Export SGF files
    - Sensei's library is an example of previous games, etc.
    - http://gokifu.com/
    - /r/baduk
    - https://smartgo.com/kifu.html
    - http://playgo.to/iwtg/en/
- C: Beginner experience
  - Possible revamp of http://playgo.to/iwtg/en/

Conclusion: let's start writing up specifications for features, so that we
can prioritize them. Plan to continue at the club on Wednesday!
