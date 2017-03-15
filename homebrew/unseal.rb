# Copyright (c) 2017 Lee Briggs
#
# This software is released under the MIT License.
#
# http://opensource.org/licenses/mit-license.php
#
class Unseal < Formula
  desc "Unseal multiple vault servers quickly and easily"
  homepage "https://github.com/jaxxstorm/unseal"
  version "v0.3-beta"

  if Hardware::CPU.is_64_bit?
    url "https://github.com/jaxxstorm/unseal/releases/download/v0.3-beta/unseal-v0.3-beta_darwin_amd64"
    sha256 "f25061683b741ad394efd4277bb076a7a1c095b98a53c76470da70aa37f08d8e"
  #else
  #  url "https://github.com/jaxxstorm/unseal/releases/download/v0.3-beta/unseal-v0.3-beta_darwin_amd64"
  #  sha256 ""
  end

  def install
    bin.install "unseal"
  end

  test do
    system "unseal version"
  end

end
