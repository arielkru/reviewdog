	"github.com/reviewdog/reviewdog/difffilter"
const diffContentAddedStrip = `diff --git a/test_added.go b/test_added.go
new file mode 100644
index 0000000..264c67e
--- /dev/null
+++ b/test_added.go
@@ -0,0 +1,3 @@
+package reviewdog
+
+var TestAdded = 14
`

			InDiff:  false,
			OldPath: "sample.old.txt",
			OldLine: 1,
			OldPath:  "sample.old.txt",
			OldLine:  0,
			InDiff:  false,
			OldPath: "nonewline.old.txt",
			OldLine: 1,
			OldPath:  "nonewline.old.txt",
			OldLine:  0,
	got := FilterCheck(results, filediffs, 0, "", difffilter.ModeAdded)
			OldPath:  "sample.old.txt",
			OldLine:  1,
			OldPath:  "sample.old.txt",
			OldLine:  0,
			OldPath:  "sample.old.txt",
			OldLine:  0,
	got := FilterCheck(results, filediffs, 0, "", difffilter.ModeDiffContext)
func TestGetOldPosition(t *testing.T) {
	const strip = 0
	tests := []struct {
		newPath     string
		newLine     int
		wantOldPath string
		wantOldLine int
	}{
		{
			newPath:     "sample.new.txt",
			newLine:     1,
			wantOldPath: "sample.old.txt",
			wantOldLine: 1,
		},
		{
			newPath:     "sample.new.txt",
			newLine:     2,
			wantOldPath: "sample.old.txt",
			wantOldLine: 0,
		},
		{
			newPath:     "sample.new.txt",
			newLine:     3,
			wantOldPath: "sample.old.txt",
			wantOldLine: 0,
		},
		{
			newPath:     "sample.new.txt",
			newLine:     14,
			wantOldPath: "sample.old.txt",
			wantOldLine: 13,
		},
		{
			newPath:     "not_found",
			newLine:     14,
			wantOldPath: "",
			wantOldLine: 0,
		},
	for _, tt := range tests {
		gotPath, gotLine := getOldPosition(filediffs, strip, tt.newPath, tt.newLine)
		if !(gotPath == tt.wantOldPath && gotLine == tt.wantOldLine) {
			t.Errorf("getOldPosition(..., %s, %d) = (%s, %d), want (%s, %d)",
				tt.newPath, tt.newLine, gotPath, gotLine, tt.wantOldPath, tt.wantOldLine)
}

func TestGetOldPosition_added(t *testing.T) {
	const strip = 1
	filediffs, _ := diff.ParseMultiFile(strings.NewReader(diffContentAddedStrip))
	gotPath, _ := getOldPosition(filediffs, strip, "test_added.go", 1)
	if gotPath != "" {
		t.Errorf("got %q as old path for addedd diff file, want empty", gotPath)