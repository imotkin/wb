pwd
echo "hello"
ps | grep "go" && echo "success" || echo "failed"
cat < example.txt | grep "e" | wc
ls -lah | grep "script" | cat