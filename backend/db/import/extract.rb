require 'csv'


tags = [
   {id: 1, name: "Ownership inscription"},
   {id: 2, name: "Gift Inscription"},
   {id: 3, name: "Date"},
   {id: 4, name: "Place name"},
   {id: 5, name: "Underlining"},
   {id: 6, name: "Annotation"},
   {id: 7, name: "Commentary"},
   {id: 8, name: "Indexing"},
   {id: 9, name: "Drawing"},
   {id: 10, name: "Poem"},
   {id: 11, name: "Quotation"},
   {id: 12, name: "Insert: paper"},
   {id: 13, name: "Insert: botanical"},
   {id: 14, name: "Insert: object"}]

csv_text = File.read("booktraces_submissions.csv")
csv_text.gsub!(/\\"/,"'")
csv = CSV.parse(csv_text, :headers => true)
cnt = 0
files_cnt = 0
tags_cnt = 0
skipped = 0
submit_sql="insert into submissions(id,upload_id,submitter_name,submitter_email,title,author,"
submit_sql << "publication_info,library,call_number,description,submitted_at,public) values "
files_sql="insert into submission_files(submission_id,filename) values "
tags_sql="insert into submission_tags(submission_id,tag_id) values "
csv.each do |row|
   out = {date_submitted: row[0]}
   data = row[1]
   found_desc = false
   lines = data.split("\n")
   dump = false
   if lines.count == 1 
      puts "SKIPPING #{row[2]}"
      # puts data 
      # puts "================================================================================"
      skipped += 1
      next
   end
   lines.each do |line|
      if found_desc 
         if line.include?("submitted-image")
            line.split('href="').each do |img|
               line.gsub(/\"/, "'")
               while true
                  idx = line.index("src='")
                  if idx.nil?
                     break
                  else 
                     # chop off all up to url and find the end quote of url
                     line = line[idx+5..-1]
                     idx2 = line.index("'")
                     if idx2 == -1 
                        abort("Missing closing quote on URL: #{line}")
                     end
                     # Grab URL and cut it out of line
                     url = line[0...idx2]
                     line = line[idx2+1..-1]
                     out[:files] << url.split("uploads/")[1]
                  end
               end
            end
         else
            out[:description] << line.gsub(/<br \/>/, "").gsub(/<p>/, "").gsub(/<\/p>/, "").gsub(/<\/strong>/, "").gsub(/^\s*/, "")
         end
      else
         if line.include?(">Title:<")
            out[:title] = line.split("</strong>")[1].gsub(/<br \/>/, "").strip
         end
         if line.include?(">Author:<")
            out[:author] = line.split("</strong>")[1].gsub(/<br \/>/, "").strip
         end
         if line.include?(">Publication date:<")
            out[:publication] = line.split("</strong>")[1].gsub(/<br \/>/, "").strip
         end
         if line.include?(">Library:<")
            out[:library] = line.split("</strong>")[1].gsub(/<br \/>/, "").strip
         end
         if line.include?(">Call number:<")
            out[:call_number] = line.split("</strong>")[1].gsub(/<br \/>/, "").strip
         end
         if line.include?(">Submitted by:<")
            out[:submitter_name] = line.split("</strong>")[1].gsub(/<br \/>/, "").strip
         end
         if line.include?(">Description:")
            found_desc = true
            out[:description] = ""
            out[:files] = []
         end
      end
   end
   out[:submitter_email] = row[3]
   out[:tags] = row[4].split(",") - ["Books", "News", "Marginalia", "Uncategorized"]
   out[:name] = row[2]
   out[:upload_id] = "LEGACY#{cnt+1}"
   print "."
   sub = "(#{cnt+1},\"#{out[:upload_id]}\",\"#{out[:submitter_name]}\",\"#{out[:submitter_email]}\",\"#{out[:title]}\","
   sub << "\"#{out[:author]}\",\"#{out[:publication_info]}\",\"#{out[:library]}\",\"#{out[:call_number]}\","
   sub << "\"#{out[:description]}\",\"#{out[:date_submitted]}\",1)"
   if cnt > 0 
      submit_sql << ","
   end
   submit_sql << sub
   
   out[:files].each do |f|
      if files_cnt > 0 
         files_sql << ","
      end
      files_sql << "(#{cnt+1},\"#{f}\")"
      files_cnt += 1
   end

   out[:tags].each do |t|
      if tags_cnt > 0 
         tags_sql << ","
      end
      tag_id = -1
      tags.each do |tag|
         if tag[:name].downcase == t.downcase
            tags_sql << "(#{cnt+1},#{tag[:id]})"
            tags_cnt += 1
         end
      end
   end

   cnt += 1
end
puts "DONE! #{cnt} records generated, #{skipped} SKIPPED; writing SQL"
f = File.open("import_wp.sql", 'w')
f.write(submit_sql)
f.write(";\n")
f.write(files_sql)
f.write(";\n")
f.write(tags_sql)
f.write(";\n")
puts "DONE!!"
