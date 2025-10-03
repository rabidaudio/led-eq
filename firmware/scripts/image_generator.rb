require 'chunky_png'
require 'optparse'
require 'nokogiri'
require 'net/http'

options = {
    width: 32,
    height: 16,
    x_offset: 0,
    y_offset: 1
}
OptionParser.new do |opts|
    opts.banner = "Usage: ruby image_generator.rb -w WIDTH -h HEIGHT -x 0 -y 1 -c COUNTRY_NAME"
    opts.on("-w", "--width") { |w| options[:width] = w }
    opts.on("-h", "--height") { |h| options[:height] = h }
    opts.on("-x", "--x-offset") { |x| options[:x_offset] = x }
    opts.on("-y", "--y-offset") { |y| options[:y_offset] = y }
    opts.on("-c", "--country COUNTRY_NAME") do |c|
        options[:country] = c
    end
end.parse!

FLAG_URL = "https://r74n.com/pixelflags/"
data = Nokogiri::HTML.parse(Net::HTTP.get(URI(FLAG_URL)))

flag_png = "#{FLAG_URL}#{data.at_css("##{options[:country]} td img").attr('src')}"

img = ChunkyPNG::Image.from_io(StringIO.new(Net::HTTP.get(URI(flag_png))))

# file output format:
# r,g,b column-wise, then row-wise
(options[:y_offset]...(options[:y_offset]+options[:height])).each do |y|
    (options[:x_offset]...(options[:x_offset]+options[:width])).each do |x|
        puts("// [#{x}, #{y}]")
        r = ChunkyPNG::Color.r(img[x, y])
        g = ChunkyPNG::Color.g(img[x,y])
        b = ChunkyPNG::Color.b(img[x, y])
        puts("%02X %02X %02X" % [r, g, b])
    end
end
