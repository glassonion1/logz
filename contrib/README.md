# logz contribution packages
Code contained in this directory contains contribution for 3rd-party Go packages.

## Contribution Packages
The following contribution packages are provided for popular Go packages and use-cases.

| Contribution Packages |
| ---- |
| [github.com/gin-gonic/gin](https://github.com/glassonion1/logz/tree/main/contrib/github.com/gin-gonic/gin/logzgin) |
| [github.com/labstack/echo](https://github.com/glassonion1/logz/tree/main/contrib/github.com/labstack/echo/logzecho) |
| [google.golang.org/grpc(alpha)](https://github.com/glassonion1/logz/tree/main/contrib/google.golang.org/grpc/logzgrpc) |

## Packaging

All contribution packages SHOULD be of the form:

```
github.com/glassonion1/logz/contrib/{IMPORT_PATH}/logz{PACKAGE_NAME}
```

Where the {IMPORT_PATH} and {PACKAGE_NAME} are the standard Go identifiers for the package being contributed.

For example:
* github.com/glassonion1/logz/contrib/github.com/gin-gonic/gin/logzgin
