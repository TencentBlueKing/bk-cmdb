Pod::Spec.new do |s|
  s.name          = "Thrift-swift3"
  s.version       = "0.12.0"
  s.summary       = "Apache Thrift is a lightweight, language-independent software stack with an associated code generation mechanism for RPC."
  s.description   = <<-DESC
The Apache Thrift software framework, for scalable cross-language services development, combines a software stack with a code generation engine to build services that work efficiently and seamlessly between C++, Java, Python, PHP, Ruby, Erlang, Perl, Haskell, C#, Cocoa, JavaScript, Node.js, Smalltalk, OCaml and Delphi and other languages.
                    DESC
  s.homepage      = "http://thrift.apache.org"
  s.license       = { :type => 'Apache License, Version 2.0', :url => 'https://www.apache.org/licenses/LICENSE-2.0' }
  s.author        = { "Apache Thrift Developers" => "dev@thrift.apache.org" }
  s.ios.deployment_target = '9.0'
  s.osx.deployment_target = '10.10'
  s.requires_arc  = true
  s.source        = { :git => "https://github.com/apache/thrift.git", :tag => "0.12.0" }
  s.source_files  = "lib/swift/Sources/*.swift"
end
