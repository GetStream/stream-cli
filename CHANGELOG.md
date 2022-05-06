# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.1.0](https://github.com/GetStream/stream-cli/compare/v1.0.0...v1.1.0) (2022-05-06)


### Features

* **channel:** extra data to create channel ([#109](https://github.com/GetStream/stream-cli/issues/109)) ([b4baeb0](https://github.com/GetStream/stream-cli/commit/b4baeb092f5a64e2f5dc7c8498e261d9c301a4f8))
* **translate:** add translate message ([#110](https://github.com/GetStream/stream-cli/issues/110)) ([cbcbd95](https://github.com/GetStream/stream-cli/commit/cbcbd95c457836bb0ae41ae53ecd2caee5bca927))


### Bug Fixes

* **docs:** add short description to deactivateuser ([0a78806](https://github.com/GetStream/stream-cli/commit/0a78806ee8d2f8fb18753d8823ee647afbf3696d))
* **upsert:** fix printing upsert user result ([b0dc374](https://github.com/GetStream/stream-cli/commit/b0dc374fb094ccdfed69b4190ceea7a422a1a7da))
* **watch_cmd:** println instead of print ([4e409f1](https://github.com/GetStream/stream-cli/commit/4e409f1585bdf13ee41b8fb230491b41ff9ceb50))

## [1.0.0](https://github.com/GetStream/stream-cli/compare/v0.3.0...v1.0.0) (2022-05-03)

### 🚨 BREAKING CHANGE 🚨
This is the first version where we rewrote the CLI from NodeJS to Go. It 
has every feature that the previous had but the interface is a bit different.
You can find the documentation in the [docs](./docs/) folder.

### Features

* add app settings crud ([d9114d7](https://github.com/GetStream/stream-cli/commit/d9114d7242e8a61904d92e2853c1e5cbd60dfeee))
* add channel crud ([5d1ebc9](https://github.com/GetStream/stream-cli/commit/5d1ebc969f9ad22e4570199e5165bab80821201a))
* add channel type crud ([3d3fee0](https://github.com/GetStream/stream-cli/commit/3d3fee04e1748a6bbf49946c4aa3c0e9ef96d9f5))
* add device crud ([045f3c8](https://github.com/GetStream/stream-cli/commit/045f3c84fd5e22f3a3e3f5bad4c12dc75ca8f0ee))
* add flag message feature ([4adcd91](https://github.com/GetStream/stream-cli/commit/4adcd919500e8bb24de4a1bb9d65068f618759fc))
* add generic viewer ([25d4d26](https://github.com/GetStream/stream-cli/commit/25d4d2662beaa07b8dcb64edcebab48b834f98db))
* add image and file upload ([fa6a5f2](https://github.com/GetStream/stream-cli/commit/fa6a5f2626808ef310a99173b369be35c051ad10))
* add import commands ([86a94b9](https://github.com/GetStream/stream-cli/commit/86a94b918a6eba1605cbc7be55b71b51a24a4502))
* add messages crud ([ed75ea8](https://github.com/GetStream/stream-cli/commit/ed75ea8abd23cdff30a01e4a9df64c30f5705951))
* add missing channel commands ([fd29250](https://github.com/GetStream/stream-cli/commit/fd292501a9c0786221477ae0cae755a25245e1db))
* add reactions crud ([#88](https://github.com/GetStream/stream-cli/issues/88)) ([a4eeab4](https://github.com/GetStream/stream-cli/commit/a4eeab45c37f0132d9b2be28383ce7121275b7e0))
* add token revocation ([#83](https://github.com/GetStream/stream-cli/issues/83)) ([3074a1d](https://github.com/GetStream/stream-cli/commit/3074a1d74ddf7c24a2c45c14c2d8b3b2642e3de5))
* add user crud ([db006b2](https://github.com/GetStream/stream-cli/commit/db006b2df9a1e1f7351e7064deceb58d421b2856))
* **push_notification:** add pushprovider update and test ([#92](https://github.com/GetStream/stream-cli/issues/92)) ([ef11aa6](https://github.com/GetStream/stream-cli/commit/ef11aa6aae78bb3444285e3c404e0cb0b5ab320a))
* set quality control ([#90](https://github.com/GetStream/stream-cli/issues/90)) ([400bb21](https://github.com/GetStream/stream-cli/commit/400bb21cd02da95375a391d521709b1b708d9977))
* switch to cobra ([d82935b](https://github.com/GetStream/stream-cli/commit/d82935bf57184789d3d6e1a05101d00d4aea8faf))
* **user:** add missing commands for users ([#97](https://github.com/GetStream/stream-cli/issues/97)) ([6944432](https://github.com/GetStream/stream-cli/commit/694443263a760bec308a4147aa615b7644e1ce72))


### Bug Fixes

* fix comments on the PR ([ed0740e](https://github.com/GetStream/stream-cli/commit/ed0740ec15315ee63c7bb81c14f8b56b504f169f))
* fix get or create channel ([e6a9657](https://github.com/GetStream/stream-cli/commit/e6a96570c03f95c07e9a0b47dee0d4737e74bed7))
* fix per pr comments ([caed0af](https://github.com/GetStream/stream-cli/commit/caed0af6edd1c2ca7ce433f82a3d56451221ad38))
* fix pr comments ([2208fab](https://github.com/GetStream/stream-cli/commit/2208fab893b8acd58e8a337f4676c01941fcdb5f))
* fix remove contenttype ([9b6a67b](https://github.com/GetStream/stream-cli/commit/9b6a67b6a07c8cc17c007781a3095cce59a4d113))
* fixes per pr comments ([4d19e00](https://github.com/GetStream/stream-cli/commit/4d19e0002c26e4653c10b2d97b4ae1169b40abba))
* fixing pr comments ([c70a5bd](https://github.com/GetStream/stream-cli/commit/c70a5bd8451c4eef762a663e7722d436050f9841))
* golangcilint fixes ([b3d82b0](https://github.com/GetStream/stream-cli/commit/b3d82b03835180559c9a67969a90a0c92947a9eb))
* pr comment fixes ([296a7ee](https://github.com/GetStream/stream-cli/commit/296a7eee5fb57c9d459f5d6bd119ceacd8c8923f))
* pr fixes ([62e66f4](https://github.com/GetStream/stream-cli/commit/62e66f4ac3407f0f1d7e511e47a8b4c53ae931ba))
* pr review fixes ([fa5afe4](https://github.com/GetStream/stream-cli/commit/fa5afe4e2261b16020a0bae32b401e14520fc15e))
* remove content type ([fc2fe38](https://github.com/GetStream/stream-cli/commit/fc2fe389e123fe6831da66e7427976c970cb5b8c))
* rename name to id ([0e66c2d](https://github.com/GetStream/stream-cli/commit/0e66c2dde58c3fa9431fed5d431ab657b569f554))

# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.
