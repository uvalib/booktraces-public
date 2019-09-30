require 'csv'


csv_text = File.read("updated_library_names.csv")
csv = CSV.parse(csv_text, :headers => true)
cnt = 0
csv.each do |row|
   puts "Convert #{row[0]} => #{row[1]}"
   cnt += 1
end
puts "DONE. Converted #{cnt} names"