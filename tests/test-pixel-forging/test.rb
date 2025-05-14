$LOAD_PATH.unshift(File.expand_path('./lib/', __dir__))
require 'grpc'
require 'pixelforging_pb'
require 'pixelforging_services_pb'

input_file = "./Logo-Pixel-Forging.png"

unless File.exist?(input_file)
  puts "Arquivo não encontrado!"
  exit
end

file_name = File.basename(input_file)
file_type = File.extname(input_file)

# Lê os bytes e divide em chunks binários
file_bytes = File.binread(input_file)
chunks = file_bytes.bytes.each_slice(32).map { |slice| slice.pack('C*') }

# Prepara a estrutura de dados no formato gRPC
requests = chunks.map do |chunk|
  PixelforgingGrpc::ExtractPaletteInput.new(
    fileBytes: chunk,
    fileName: file_name, 
    fileType: file_type
  )
end
outputbytes = []

stub = PixelforgingGrpc::PixelForging::Stub.new('localhost:9090', :this_channel_is_insecure) 
stub.extract_palette(requests.each) do  |r|
  outputbytes <<  r.paletteBytes
end
File.open("output2.png", "wb") do |f|
  outputbytes.each do |chunk|
    f.write(chunk)
  end
end
puts "PNG recebido e salvo como 'output.png'"
