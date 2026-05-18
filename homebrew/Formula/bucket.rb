# Draft formula for the bucket tap.
# The tap is expected to be added with:
#   brew tap jpwain/bucket https://github.com/jpwain/bucket.git

class Bucket < Formula
  desc "Terminal tool for reviewing two text files side by side"
  homepage "https://github.com/jpwain/bucket"
  version "0.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_darwin_arm64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_FOR_DARWIN_ARM64"
    else
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_darwin_amd64.tar.gz"
      sha256 "REPLACE_WITH_SHA256_FOR_DARWIN_AMD64"
    end
  end

  def install
    bin.install "bucket"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/bucket --version")
  end
end
