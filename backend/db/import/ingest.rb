require 'roo'
require 'fileutils'

# copy images from src into proper structure at dest. return insert sql
def process_images(sub_id, images, src_dir, dest_dir)
   sql = ""
   rel_path = dest_dir.split("submitted/")[1]
   images.each do |i|
      puts "IMG: [#{i}]"
      # make sure dest dir exists and copy src file 
      dest_fn = File.join(dest_dir, i)
      dest_path = File.dirname dest_fn
      if Dir.exist?(dest_path) == false
         FileUtils.mkdir_p(dest_path)
      end
      src_fn = File.join(src_dir, i) 
      if File.exist?(src_fn) == false 
         abort "ERROR: #{src_fn} not found. Skipping"
         next
      end
      # puts "CP #{src_fn} -> #{dest_fn}"
      FileUtils.cp(src_fn, dest_fn)

      if sql != "" 
         sql += ","
      end
      fn = File.join(rel_path, i)
      sql +=  "(#{sub_id},\"#{fn}\")"
   end
   return sql
 end

 ##
 ## MAIN ##
 ##   ruby ingest.rb mary_baldwin_bt.xlsx "Mary Baldwin" ~/Desktop/mary_baldwin_img ~/dev/booktraces-dev/booktraces-public/submissions/submitted/2019/ 1054
 ##
fn = ARGV[0]
univ = ARGV[1]
img_dir = ARGV[2]
dest_dir = ARGV[3]
id = ARGV[4].to_i
abort "Required: xlsx univ_name image_dir dest_dir start_id" if fn.empty? || 
   univ.empty? || img_dir.empty? || dest_dir.empty? || id == 0

univ_fn = univ.downcase.gsub(/\s/, "_")
dest_dir = File.join(dest_dir, "3cavaliers", univ_fn)
out = univ_fn + ".sql"

sub_name = "Book Traces: Civil War Readers and Their Books in Virginia Libraries"
sub_email = "khj5c@virginia.edu"
submit_sql="insert into submissions(id,upload_id,submitter_name,submitter_email,title,author,"
submit_sql << "publication_info,library,call_number,description,submitted_at,public) values "
files_sql="insert into submission_files(submission_id,filename) values "

puts "Ingest #{univ}:  #{fn}, images: #{img_dir}\n\tDestination: #{dest_dir}, StartID: #{id}"
puts "========================================================================================"
xlsx = Roo::Excelx.new(fn)
img_cnt = 0
xlsx.each_with_pagename do |name, sheet|
   puts "Process  sheet #{name}"
   sheet.parse(title: 'Title', author: "Author" , library: /Location*/, 
      call_num: "Call Number", desc: "Notes",
      image1: "Image_1", image2: "Image_2", image3: "Image_3",clean:true).each do |row|
      
      images = []
      images << row[:image1] if !row[:image1].nil? && !row[:image1].empty?
      images << row[:image2] if !row[:image2].nil? && !row[:image2].empty?
      images << row[:image3] if !row[:image3].nil? && !row[:image3].empty?
      if !images.empty?
         puts "     #{row[:call_num]} has interventions"
         # generate submission insert
         upload_id="3cavaliers#{id}"
         lib = "#{row[:library]}, #{univ}"
         sub = "(#{id},\"#{upload_id}\",\"#{sub_name}\",\"#{sub_email}\",\"#{row[:title]}\","
         sub << "\"#{row[:author]}\",\"\",\"#{lib}\",\"#{row[:call_num]}\","
         sub << "\"#{row[:desc]}\",\"#{DateTime.now()}\",1)"
         if img_cnt > 0 
            submit_sql << ",\n"
         end
         submit_sql << sub

         # move images to destination and generate SQL
         img_sql = process_images(id, images, img_dir, dest_dir)
         if img_cnt > 0 
            files_sql += ",\n"
         end
         files_sql += img_sql
         img_cnt+=1
         id+=1
      end
   end
 end

 puts "DONE! #{img_cnt} records generated, writing SQL"
f = File.open("import_#{univ_fn}.sql", 'w')
f.write(submit_sql)
f.write(";\n")
f.write(files_sql)
f.write(";\n")
puts "DONE!!"