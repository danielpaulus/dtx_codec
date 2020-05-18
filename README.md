# dtx_codec
A golang based Apple DTX implementation. So you can run on real iOS devices: XCUITests, get CPU metrics, launch and kill apps from any OS without the need for expensive Apple hardware :-)

Will be added to go-ios eventually.
Use https://github.com/danielpaulus/ios_simulator_dtx_dump to get a dump of DTX messages to test the decoder with.

Done:
- Basic Decoder, fully decoding DTX messages and dump them

Check out this example method call, which the device sends to us to tell us about the `blaUITests.blaUITests` testcase finishing:
```
i1038.0e c1 t:rpc_asking_reply mlen:25357 aux_len25162 paylen179
auxheader:BufSiz:25584 Unknown:0 AuxSiz:25146 Unknown2:0
aux:[{t:binary, v:["blaUITests.blaUITests"]},
{t:binary, v:["blaUITests.blaUITests"]},
{t:binary, v:["blaUITests.blaUITests"]},
]
payload: "_XCT_testCase:method:didFinishActivity:" 
```
 
 Todo:
- Basic Encoder, re-encode DTX so you can control stuff
- Fix a few unknown things for real devices (I am using Simulator output to develop before switching to devices)
- Integrate into go-ios

Check out my nskeyedarchiver implementation here: https://github.com/danielpaulus/nskeyedarchiver
