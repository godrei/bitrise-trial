# IPA analyzer summary

## Understanding the requirements

Based on the description I identified the following steps:

1. Getting the contents: since .ipa files are archives with a specific file structure, the first step was to unarchive them.
2. Finding the relevant information: I wanted to use the archive's `Info.plist` file to get the necessary values. Based on [Apple's documentation](https://developer.apple.com/library/archive/documentation/General/Reference/InfoPlistKeyReference/Articles/CoreFoundationKeys.html#//apple_ref/doc/uid/TP40009249-SW10) I used the following keys:
    - Bundle Identifier: `CFBundleIdentifier`
    - Version number: `CFBundleShortVersionString`
    - Build number: `CFBundleVersion`
    - Icon files: `CFBundleIconFiles`. 
    
    I wasn't entirely sure about `CFBundleIconFiles`, but the following segment of the [relevant docs](https://developer.apple.com/library/archive/documentation/CoreFoundation/Conceptual/CFBundles/BundleTypes/BundleTypes.html#//apple_ref/doc/uid/10000123i-CH101-SW1) convinced me:
    > How you identify these images to the system can vary, but the recommended way to specify your application icons is to use the CFBundleIconFiles key.

    My initial idea was to use either `plutil`, or `defaults` to read these values

3. Printing the values.

Based on these steps I decided to quickly write a shell script:

```bash
#!/bin/bash

IPA_PATH=$1
UNZIPPED_IPA_PATH="contents"
IPA_PAYLOAD="${UNZIPPED_IPA_PATH}/Payload"

unzip "${IPA_PATH}" -d "${UNZIPPED_IPA_PATH}"
APP_NAME=$(ls $IPA_PAYLOAD)

PLIST_PATH="$(pwd)/${IPA_PAYLOAD}/${APP_NAME}/Info.plist"
echo "Bundle Identifier: $(defaults read $PLIST_PATH CFBundleIdentifier)"
echo "Version number: $(defaults read $PLIST_PATH CFBundleShortVersionString)"
echo "Build number: $(defaults read $PLIST_PATH CFBundleVersion)"
echo "Icon files: $(defaults read $PLIST_PATH CFBundleIconFiles)"
```

## Translating the idea into Go

The next challenge was to figure out how to do these steps in Go. In the first approximation I decided to keep using `defaults` to read the values:

```go
func readPlistValue(key string, plistPath string) string {
    command := exec.Command("defaults", "read", plistPath, key)
    out, err := command.CombinedOutput()

    if err != nil {
        log.Fatal(err)
    }

    return string(out)
}
```

What bugged me about my first solution was that it relied on `plutil` or `defaults` to do the heavy lifting, so it assumed a macOS environment. If users wanted to run the analyzer on the same machine they use for creating the IPA archive, this assumption makes some sense but I wanted to see if I have any other options for doing the parsing.

## Challenges

I think it's safe to say that the biggest roadblock was the binary plist format. While I found some good [resources](https://medium.com/@karaiskc/understanding-apples-binary-property-list-format-281e6da00dbd) explaining the format, I felt that the combined challenge of understanding it, and implementing a parser in a language I have no experience is a bit more than I can chew for the available time. I started looking for third-party packages, and found [go-plist](https://godoc.org/github.com/DHowett/go-plist), which proved to be the right tool for the job.

I wasn't sure how to package `go-plist` together with my code, so I decided to run `go build` to produce the `bitrise-trial` executable, and include it with the package.

Another challenging part was getting the plist values in the correct type. Since `go-plist` uses the `map[string]interface{}` type for plist dictionaries, I had to use some reflection to convert the values into the desired types. This was especially tricky with the icon files.

## Testing

As a final step I decided to write a few tests for the plist parsing logic. Although the heavy lifting is done by the `go-plist` package, I wrote the unmarshalling code, and I wanted to ensure that the values are being mapped to the `plistInfo` type.

## Timing

I tried to keep track of the time I spent doing this exercise. These are not exact values, but rough estimates:
- Hacky shell script implementation: 30 minutes, mostly spent by reading the docs
- Go implementation with `defaults`: 1:30 hours. Most of it was experimenting with the language, and figuring out basic stuff, like reading files, and unarchiving.
- Reading about the binary plist format: 1 hour. Since I opted for a third-party library to do the parsing, this wasn't necessarily time well spent :)
- Integraging `go-plist`: 1-1:30 hours. Figuring out how to properly implement the `UnmarshalPlist` method took some time. I also had to fix my reflection-handling code after a failing test. :)
- Writing tests: around 30 minutes.
