require 'csv'


csv_text = File.read("booktraces_news.csv")
csv_text.gsub!(/\\"/,"'")
csv = CSV.parse(csv_text, :headers => true)
cnt = 0
news_sql="insert into news(id,title,content,published,created_at) values "
csv.each do |row|
   out = {date_submitted: row[0]}

   print "."
   news = "(#{cnt+1},\"#{row[2]}\",\"#{row[1]}\",1,\"#{out[:date_submitted]}\")"
   if cnt > 0 
      news_sql << ","
   end
   news_sql << news
   

   cnt += 1
end
puts "DONE! #{cnt} records generated"
f = File.open("import_news.sql", 'w')
f.write(news_sql)
f.write(";\n")