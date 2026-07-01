const parser = require("@babel/parser");
const fs = require("fs");

const filePath = process.argv[2];
const code = fs.readFileSync(filePath, "utf8");

try {
  const ast = parser.parse(code, {
    sourceType: "unambiguous",
    errorRecovery: true,
    plugins: [
      "jsx",
      "typescript",
      "classProperties"
    ]
  });

  console.log(JSON.stringify(ast, null, 2));
} catch (err) {
  console.error("Parse error:", err.message);
  process.exit(1);
}

