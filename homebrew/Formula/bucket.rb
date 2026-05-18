# Draft formula for the bucket tap.
# The tap is expected to be added with:
#   brew tap jpwain/bucket https://github.com/jpwain/bucket.git

class Bucket < Formula
  desc "Terminal tool for reviewing two text files side by side"
  homepage "https://github.com/jpwain/bucket"
  version "1.1.0"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_darwin_arm64.tar.gz"
      sha256 "89af4d96ace8efa384034aaf4fe072a67f3b5ac568aa9baefd34d608326e9c86"
    else
      url "https://github.com/jpwain/bucket/releases/download/v#{version}/bucket_#{version}_darwin_amd64.tar.gz"
      sha256 "e8ae44ff3ca073aeee060bfb19cb6717d1cadb91ad9f85d518e35af390d3de48"
    end
  end

  def install
    bin.install "bucket"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/bucket --version")
  end
end
