# streamdeck-LMS-intermediary

Simple example of providing HTTP endpoints via a locally-run server which are then used to POST data to [Logitech Media Server](https://en.wikipedia.org/wiki/Logitech_Media_Server) (LMS).

These HTTP endpoints are then able to be called from the [Elgato Streamdeck](https://en.wikipedia.org/wiki/Elgato#Stream_Deck) such that the Elgato Streamdeck can be used to control the LMS.

For more context, please see this blog post: <https://algrt.hm/2023-06-29-elgato-streamdeck/>.

To build, simply run 

```shell
go build .
```

in this directory, on a machine which has the go toolchain available.