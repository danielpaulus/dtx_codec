# dtx_codec
A golang based Apple DTX implementation. So you can run on real iOS devices: XCUITests, get CPU metrics, launch and kill apps from any OS without the need for expensive Apple hardware :-)

Will be added to go-ios eventually.
Use https://github.com/danielpaulus/ios_simulator_dtx_dump to get a dump of DTX messages to test the decoder with.

Done:
- Basic Decoder, fully decoding DTX messages and dump them
 
 Todo:
- Basic Encoder, re-encode DTX so you can control stuff
- Fix a few unknown things for real devices (I am using Simulator output to develop before switching to devices)
- Integrate into go-ios
Check it out here: https://github.com/danielpaulus/nskeyedarchiver
