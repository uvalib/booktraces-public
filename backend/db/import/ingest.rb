require 'roo'
require 'fileutils'
require 'find'

$img_cnt = 0
$err_cnt = 0

# copy images from src into proper structure at dest. return insert sql
def process_images(sub_id, images, src_dir, process, dest_dir)
   sql = ""
   rel_path = dest_dir.split("submitted/")[1]

   # Note: i in this loop is the relative path to the image;
   # for example - B_Images/JMU-BX-0005.Image_2.033932.jpg
   images.each do |i|
      puts "  IMG: [#{i}] REL #{rel_path} - Process: #{process}"

      
      # Look for source image in correct matching location
      src_fn = File.join(src_dir, i) 
      if File.exist?(src_fn) == false 
         # Not found; just look for the filename anywhere in the source tree
         img_fn = "**/#{File.basename(i)}"
         hits = Dir.glob(File.join(src_dir, img_fn))
         if hits.length == 0 
            $err_cnt += 1
            puts "ERROR: #{img_fn} not found. Skipping"
            next
         elsif hits.length > 1 
            $err_cnt += 1
            puts "ERROR: Multiple hits for #{img_fn} found. Skipping"
            next
         else 
            src_fn = hits.first
         end
      end
      
      if process == true
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
         if File.exist?(dest_fn) == false
            FileUtils.cp(src_fn, dest_fn)
         end
         if File.exist?(thumb_fn) == false 
            cmd = "convert -quiet -resize 150x150^ -extent 150x150 -gravity center \"#{dest_fn}\" \"#{thumb_fn}\""
            `#{cmd}`
         end
      else
         $img_cnt += 1
      end

      if sql != "" 
         sql += ",\n"
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
univ_id = ARGV[1]
start_id = ARGV[2]
dest_dir = ARGV[4]
process_images = ARGV[3]
if process_images == "1" || process_images == "true"
   process_images = true 
else
   process_images = false
end

puts "Ingest initializing..."
abort "start_id reqquired" if start_id.nil?
puts "  Starting ingest ID #{start_id}"
puts("  University ID #{univ_id}")


ingest_base = File.dirname(fn)
img_dir = File.join(ingest_base, "B_Images")
if !File.exist?(fn) 
   abort "#{fn} not found"
end
if !Dir.exist?(img_dir)
   abort "Images directory #{img_dir} not found"
end
univ_fn = File.basename(fn).split(".")[0]
puts "  Ingest base name: #{univ_fn}"
puts "  Ingest: #{fn}\n  Images: #{img_dir}"
puts "  Process Images? #{process_images}"
abort "dest_dir is required" if process_images && dest_dir.nil? 

if !Dir.exist?(dest_dir)
   abort "Image Destination directory #{dest_dir} not found"
end
dest_dir = File.join(dest_dir, "partners", univ_fn)
puts "  Images Dest Dir: #{dest_dir}"  

sub_name = "Book Traces: Partner Library Survey"
sub_email = "khj5c@virginia.edu"
submit_sql="insert into submissions(id,upload_id,submitter_name,submitter_email,title,author,"
submit_sql << "publication_info,library,call_number,description,submitted_at,public,institution_id) values \n"
files_sql = "insert into submission_files(submission_id,filename) values \n"

puts
puts "=================="
puts "STARTING INGESTION"
puts "=================="
xlsx = Roo::Excelx.new(fn)
intervention_cnt = 0
tgt_sheets = [ "all one table"]
id = start_id.to_i
skips = 0
xlsx.each_with_pagename do |name, sheet|
   if tgt_sheets.include?(name) == false 
      next
   end
   puts "Process #{name}"
   rows = []
   if univ_id.to_i != 77
      rows = sheet.parse(title: 'Title', author: "Author" , library: "Location-Building", 
         call_num: "Call Number", desc: "Notes", image_bc: "Image_barcode",
         interventions: "Interventions",
         image1: "Image_1", image2: "Image_2", image3: "Image_3", clean:true)
   else
      rows = sheet.parse(title: 'Title', author: "Author" , library: "Location-Building", 
         call_num: "Call number concatenate", desc: "Notes", interventions: "Interventions",
         image1: "Image_1", image2: "Image_2", image3: "Image_3", clean:true)
   end
   
   rows.each do |row|
      images = []
      if univ_id.to_i != 77
         images << row[:image_bc] if !row[:image_bc].nil? && !row[:image_bc].empty? && !row[:image_bc].include?("Unable")
      end
      images << row[:image1] if !row[:image1].nil? && !row[:image1].empty? && !row[:image1].include?("Unable")
      images << row[:image2] if !row[:image2].nil? && !row[:image2].empty? && !row[:image2].include?("Unable")
      images << row[:image3] if !row[:image3].nil? && !row[:image3].empty? && !row[:image3].include?("Unable")
      if !images.empty? && row[:interventions] != "No interventions" && row[:interventions] != "Not on shelf"
         puts "     #{row[:call_num]} has interventions (curr count: #{intervention_cnt})"
         # generate submission insert
         upload_id="#{univ_fn}#{intervention_cnt+1}"
         lib = "#{row[:library]}"
         ingest_date = DateTime.now().strftime("%F %T")
         title = row[:title]
         if title.index('/') == title.length-1
            title = title[0..-2].strip 
         end
         title = title.gsub(/\"/, '\\"')
         desc = row[:desc]
         if !desc.nil?
            desc = desc.gsub(/\"/, '\\"')
         end

         # move images to destination and generate SQL
         img_sql = process_images(id, images, img_dir, process_images, dest_dir)
         if img_sql.empty? 
            puts "ERROR: No images found for interventions on #{row[:call_num]}"
            skips+= 1
         else
            sub = "(#{id},\"#{upload_id}\",\"#{sub_name}\",\"#{sub_email}\",\"#{title}\","
            sub << "\"#{row[:author]}\",\"\",\"#{lib}\",\"#{row[:call_num]}\","
            sub << "\"#{desc}\",\"#{ingest_date}\",1,#{univ_id})"
            if intervention_cnt > 0 
               submit_sql << ",\n"
            end
            submit_sql << sub

            if intervention_cnt > 0 
               files_sql += ",\n"
            end
            files_sql += img_sql

            intervention_cnt+=1
            id+=1
         end
      end
   end
 end

sql_out = "import_#{univ_fn}.sql"
puts "DONE! #{intervention_cnt} (#{skips} skipped) records generated with #{$img_cnt} images and #{$err_cnt} image errors, writing SQL to #{sql_out}"
if File.exist?(sql_out) 
   puts "  * Removing #{sql_out}"
   File.delete(sql_out)
end
f = File.open(sql_out, 'w')
f.write(submit_sql)
f.write(";\n")
f.write(files_sql)
f.write(";\n")
puts "DONE!!"


# ruby ingest.rb ~/Desktop/bt_ingest_2020/brynmawr/bryn_mawr.xlsx  73 2461 0  ~/dev/booktraces-dev/booktraces-public/submissions/submitted/2020 
# ruby ingest.rb ~/Desktop/bt_ingest_2020/uconn/uconn.xlsx  44 2627 0  ~/dev/booktraces-dev/booktraces-public/submissions/submitted/2020 
# ruby ingest.rb ~/Desktop/bt_ingest_2020/brandeis/brandeis.xlsx  77 2849 0  ~/dev/booktraces-dev/booktraces-public/submissions/submitted/2020 

