Using ColorLight 5A-75B board to drive the display. It works over ethernet.
Unfortunately:
- configuring the display requires the [LEDVISION](https://en.colorlightinside.com/product/download/381) software which is windows only
- The protocol is a bit opaque, but [some folks have reverse-engineered it](https://hkubota.wordpress.com/2022/01/31/winter-project-colorlight-5a-75b-protocol/)

[This code](https://github.com/haraldkubota/colorlight) has an implementation in dart, but it relies on a linux-only ethernet driver.

