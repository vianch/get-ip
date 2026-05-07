class GetIp < Formula
  desc "Beautiful TUI showing your network interfaces, IPs, MACs, gateways, and DHCP/Manual mode"
  homepage "https://github.com/vianch/get-ip"
  url "https://github.com/vianch/get-ip/archive/refs/tags/v0.1.0.tar.gz"
  sha256 "REPLACE_WITH_RELEASE_TARBALL_SHA256"
  license "MIT"
  head "https://github.com/vianch/get-ip.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X main.version=#{version}
    ]
    system "go", "build", *std_go_args(ldflags: ldflags), "./"
  end

  test do
    assert_match "get-ip", shell_output("#{bin}/get-ip --version")
  end
end
