# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class StreamCli < Formula
  desc "Manage your Stream applications easily."
  homepage "https://github.com/GetStream/stream-cli"
  version "1.7.2"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/GetStream/stream-cli/releases/download/v1.7.2/stream-cli_Darwin_arm64.tar.gz"
      sha256 "731abfc2a08440861e9340e5989b8a5357677bb2fc352aa6cae24d5e9b8c54e3"

      def install
        bin.install "stream-cli"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/GetStream/stream-cli/releases/download/v1.7.2/stream-cli_Darwin_x86_64.tar.gz"
      sha256 "90d2442c6071647ca5a20b28e373e4ee030b664f10e589b6355e2a245dc27d46"

      def install
        bin.install "stream-cli"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/GetStream/stream-cli/releases/download/v1.7.2/stream-cli_Linux_arm64.tar.gz"
      sha256 "5608f00b633c1ccec0ee65cadd42a9765337838269024c3f23003ceaa2feaf40"

      def install
        bin.install "stream-cli"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/GetStream/stream-cli/releases/download/v1.7.2/stream-cli_Linux_x86_64.tar.gz"
      sha256 "274fd30b6635d40135390f4cf96c33ca8770a230718d460ffe0ad10eaa62ed02"

      def install
        bin.install "stream-cli"
      end
    end
  end
end
