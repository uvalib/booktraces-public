require 'roo'
require 'fileutils'
require 'find'

$img_cnt = 0

# copy images from src into proper structure at dest. return insert sql
def process_images(sub_id, images, src_dir, dest_dir)
   sql = ""
   rel_path = dest_dir.split("submitted/")[1]

   # Note: i in this loop is the relative path to the image;
   # for example - B_Images/JMU-BX-0005.Image_2.033932.jpg
   images.each do |i|
      puts "IMG: [#{i}]"

      # Look for source image in correct matching location
      src_fn = File.join(src_dir, i) 
      if File.exist?(src_fn) == false 
         # Not found; just look for the filename anywhere in the source tree
         img_fn = "**/#{File.basename(i)}"
         hits = Dir.glob(File.join(src_dir, img_fn))
         if hits.length == 0 
            puts "ERROR: #{img_fn} not found. Skipping"
            next
         elsif hits.length > 1 
            puts "ERROR: Multiple hits for #{img_fn} found. Skipping"
            next
         else 
            src_fn = hits.first
         end
      end
   
      # make sure dest dir exists and copy src file 
      dest_fn = File.join(dest_dir, i)
      dest_path = File.dirname dest_fn
      if Dir.exist?(dest_path) == false
         FileUtils.mkdir_p(dest_path)
      end
      
      # puts "CP #{src_fn} -> #{dest_fn}" 
      $img_cnt += 1
      base_fn = File.basename(dest_fn, File.extname(dest_fn))
      tfn = "#{base_fn}-150x150#{File.extname(dest_fn).downcase}"
      thumb_fn = File.join(File.dirname(dest_fn), tfn)
      if File.exist?(dest_fn) == false &&  File.exist?(thumb_fn)
         FileUtils.cp(src_fn, dest_fn)
         cmd = "convert -quiet -resize 150x150^ -extent 150x150 -gravity center \"#{dest_fn}\" \"#{thumb_fn}\""
         `#{cmd}`
      end

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
intervention_cnt = 0
tgt_sheets = [ "B", "D", "E", "PR", "PS", "U", "multivolumes", "all", "EMU list"]
xlsx.each_with_pagename do |name, sheet|
   if tgt_sheets.include?(name) == false 
      puts "skipping unrecognized sheet #{name}"
      next
   end
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
         ingest_date = DateTime.now().strftime("%F %T")
         title = row[:title].gsub(/\"/, '\\"')
         sub = "(#{id},\"#{upload_id}\",\"#{sub_name}\",\"#{sub_email}\",\"#{title}\","
         sub << "\"#{row[:author]}\",\"\",\"#{lib}\",\"#{row[:call_num]}\","
         sub << "\"#{row[:desc]}\",\"#{ingest_date}\",1)"
         if intervention_cnt > 0 
            submit_sql << ",\n"
         end
         submit_sql << sub

         # move images to destination and generate SQL
         img_sql = process_images(id, images, img_dir, dest_dir)
         if intervention_cnt > 0 
            files_sql += ",\n"
         end
         files_sql += img_sql
         intervention_cnt+=1
         id+=1
      end
   end
 end

sql_out = "import_#{univ_fn}.sql"
f = File.open(sql_out, 'w')
puts "DONE! #{intervention_cnt} records generated with #{$img_cnt} images, writing SQL to #{sql_out}"
f.write(submit_sql)
f.write(";\n")
f.write(files_sql)
f.write(";\n")
puts "DONE!!"

## UPDATES
## ALTER TABLE submissions MODIFY title TEXT NOT NULL;
## ALTER TABLE submissions MODIFY call_number varchar(100);

# sequence (XX has been ingested on prod)
# =======================================
# XX mary_baldwin (1054-1071) xxx ingested 
# XX jmu (1072-1105)
# XX VCU (1106-1110)
# XX washington_and_lee (1111-1232)
# XX emory and  henry (1233-1268)
# XX virginia_tech (1269-1301)
# XX roanoke (1302-1347)
# XX VMI (1348-1502)
# XX VUU (1503-1619)
# XX EMU (1620-1714)
# XX U Richmond (1715-1783)
# william and mary (1784-)

## ruby ingest.rb ~/Desktop/urichmond/richmond.xlsx "University of Richmond" ~/Desktop/urichmond/ ~/dev/booktraces-dev/booktraces-public/submissions/submitted/2019/ 1715
