# Changelog

## [0.1.4](https://github.com/abemedia/go-don/compare/v0.1.3...v0.1.4) (2023-05-09)


### Bug Fixes

* error parsing GET request ([#108](https://github.com/abemedia/go-don/issues/108)) ([8298a2c](https://github.com/abemedia/go-don/commit/8298a2c7a8d46858420fccbbe39909db71838b38))

## [0.1.3](https://github.com/abemedia/go-don/compare/v0.1.2...v0.1.3) (2023-05-08)


### Features

* api serve ([#98](https://github.com/abemedia/go-don/issues/98)) ([602b24c](https://github.com/abemedia/go-don/commit/602b24c5220bee9955d30ec38e7fbc8b41aa2e10))
* implement errors interfaces & improve tests ([#91](https://github.com/abemedia/go-don/issues/91)) ([0a282f0](https://github.com/abemedia/go-don/commit/0a282f0fc2fbe289a89fd9cc0ba94939108fb205))


### Bug Fixes

* decoder panics on tag not found ([#100](https://github.com/abemedia/go-don/issues/100)) ([3a73c35](https://github.com/abemedia/go-don/commit/3a73c35dd996e1035360733d4b60d52b88c3243b))


### Performance Improvements

* improve encoding performance ([#104](https://github.com/abemedia/go-don/issues/104)) ([9dbebfa](https://github.com/abemedia/go-don/commit/9dbebfa81db3277efd964d6d8fa9f1755ef9683a))

## [0.1.2](https://github.com/abemedia/go-don/compare/v0.1.1...v0.1.2) (2023-05-05)


### Features

* improve error handling ([#85](https://github.com/abemedia/go-don/issues/85)) ([3f976fc](https://github.com/abemedia/go-don/commit/3f976fca67e518b9c786c4af32c46586fd5cdc06))


### Bug Fixes

* panic on interface request ([#82](https://github.com/abemedia/go-don/issues/82)) ([8e83bd6](https://github.com/abemedia/go-don/commit/8e83bd692db5569b36426b112d4d243cc106968a))

## [0.1.1](https://github.com/abemedia/go-don/compare/v0.1.0...v0.1.1) (2023-01-14)


### Features

* better error handling, minor refactor ([#58](https://github.com/abemedia/go-don/issues/58)) ([0de3fc3](https://github.com/abemedia/go-don/commit/0de3fc32deb4692a7e768f1f650122b664785810))
* **encoding/text:** support marshaler & stringer, improve performance ([#60](https://github.com/abemedia/go-don/issues/60)) ([b7bffe8](https://github.com/abemedia/go-don/commit/b7bffe81d2ca0651a78d694462e6684df211f0ca))

## 0.1.0 (2022-08-29)


### Bug Fixes

* **encoding/text:** int conversion on 32bit ([#32](https://github.com/abemedia/go-don/issues/32)) ([3e469fe](https://github.com/abemedia/go-don/commit/3e469fe24189849d25e24395500eca23d6043a96))
* use error's marshaler, circular ref in Handler, lint issues ([#28](https://github.com/abemedia/go-don/issues/28)) ([f8b32ea](https://github.com/abemedia/go-don/commit/f8b32eaa0150d96a6ce186f2bdf41ef0e90a39e0))
