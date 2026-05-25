# Homebrew tap formula for bucket.
# This repository can be tapped directly with:
#   brew tap jpwain/bucket https://github.com/jpwain/bucket.git

class Bucket < Formula
  desc "Terminal tool for reviewing two text files side by side"
  homepage "https://github.com/jpwain/bucket"
  version "1.2.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_darwin_arm64.tar.gz"
      sha256 "52450d3b3af19b69b0ce97aca65f761ec59b12a97b8270e86d4a7512a1ac4a51"
    else
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_darwin_amd64.tar.gz"
      sha256 "b7df067d2e97e28ab4158a389534b167919cbc08249c6621670a960eda129850"
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_linux_arm64.tar.gz"
      sha256 "8beb026e8316145cdb72bb57e8431a16457c9930d3b60a2f4adfec8d2a396408"
    else
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_linux_amd64.tar.gz"
      sha256 "ad515d8a02592990abcf96e301b040cd2354b4ebe091a0f38f80334affdaa1ff"
    end
  end

  def install
    bin.install "bucket"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/bucket --version")
  end
end
