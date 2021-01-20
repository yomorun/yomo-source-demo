# yomo-source-example

The demo of [yomo-source](https://yomo.run/source) which sends random float32 data to [yomo-zipper](https://yomo.run/zipper) in every 100 ms.

## How to run the demo project

> **Note**: you need to run [yomo-zipper](https://yomo.run/zipper#how-to-config-and-run-yomo-zipper) first before running this demo project.

### Run `yomo-zipper`

In order to experience the data processing in YoMo, you can run `yomo wf run workflow.yaml` with the current example yomo-source-example as `yomo-source`. See [yomo-zipper](https://yomo.run/zipper#how-to-config-and-run-yomo-zipper) for details.

### Run `yomo-source-example`

``` shell
# YOMO_ZIPPER_ENDPOINT is the address of `yomo-zipper`.
YOMO_ZIPPER_ENDPOINT=localhost:9999 go run main.go
```

You will see the following message:

```shell
2020/12/29 15:59:21 ✅ Emit 83.809280 to yomo-zipper
2020/12/29 15:59:21 ✅ Emit 135.370453 to yomo-zipper
2020/12/29 15:59:21 ✅ Emit 180.532379 to yomo-zipper
2020/12/29 15:59:21 ✅ Emit 7.493614 to yomo-zipper
2020/12/29 15:59:21 ✅ Emit 159.445312 to yomo-zipper
2020/12/29 15:59:21 ✅ Emit 10.719324 to yomo-zipper
```

## How to write your `yomo-source`

1. Connect to `yomo-zipper` over **QUIC** and create a QUIC stream.

```go
// connect to yomo-zipper via QUIC.
client, err := quic.NewClient(zipper)
if err != nil {
  return err
}

// create a stream
stream, err := client.CreateStream(context.Background())
if err != nil {
  return err
}
```

2. Encode your data via [y3-codec](https://github.com/yomorun/y3-codec-golang).

```go
protoCodec := codes.NewProtoCodec(0x10)
sendingBuf, _ := protoCodec.Marshal(randData)
```

3. Send data to `yomo-zipper` via **QUIC stream**.

```go
_, err := stream.Write(sendingBuf)
```

## How `yomo-source` and `yomo-zipper` work

![YoMo](https://github.com/yomorun/yomo-source-demo/blob/main/yomo.png)
