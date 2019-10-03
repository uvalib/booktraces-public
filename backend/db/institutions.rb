require 'csv'


csv_text = File.read("updated_library_names.csv")
csv = CSV.parse(csv_text, :headers => true)
cnt = 0
out = ""
csv.each do |row|
   puts "Convert #{row[0]} => #{row[1]}"
   out << "update submissions set institution_id=(select id from institutions where name=\"#{row[1]}\")"
   out << " where library=\"#{row[0]}\";\n"
   cnt += 1
end
puts "DONE! #{cnt} updates generated\n\n"
puts "Generate delete of bad items"

bad_libs = ["Inquiring Minds Bookstore", "Sub basement"]
bad_libs.each do |lib|
   out << "delete from submission_tags where submission_id=(select id from submissions where library=\"#{lib}\");\n"
   out << "delete from submission_files where submission_id=(select id from submissions where library=\"#{lib}\");\n"
   out << "delete from submissions where library=\"#{lib}\";\n"
end

f = File.open("update_library.sql", 'w')
f.write(out)