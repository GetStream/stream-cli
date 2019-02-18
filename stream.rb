class Stream < Formula
  desc "The Stream Command Line Interface (CLI) makes it easy to create and manage your Stream apps directly from the terminal."
  homepage "https://getstream.io"
  url "https://stream-cli.s3.amazonaws.com/channels/beta/stream.zip"
  version "0.0.1-beta.5"
  sha256 "bdea715efaee2c4a53cfd4e0bd8726210d7475c66163dca23932283370bbf00c"

  def install
    system "./configure", "--disable-debug",
                          "--disable-dependency-tracking",
                          "--disable-silent-rules",
                          "--prefix=#{prefix}"
    system "make", "install"
  end
end
