# memos首个自动化测试

### 1. 初始化GitHub Actions
在2022年8月27日的0.4.1版本中，memos项目引入了自动化测试

在项目的根目录下创建一个.github/workflows文件夹，并在其中添加工作流文件（YAML格式），这样GitHub Actions就会自动运行这些文件。

### 2. 添加基础的测试工作流
```shell
name: Go Tests

on:
  push:
    branches:
      - main
  pull_request:
    branches:
      - main

jobs:
  go-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v3
        with:
          go-version: 1.18
          cache: true
      - name: Run Tests
        run: go test -v ./... | tee test.log; exit ${PIPESTATUS[0]}
      - name: Analyze Test Time
        run: grep --color=never -e '--- PASS:' -e '--- FAIL:' test.log | sed 's/[:()]//g' | awk '{print $2,$3,$4}' | sort -t' ' -nk3 -r | awk '{sum += $3; print $1,$2,$3,sum"s"}'
```

### 3. 推送到GitHub并观察Actions运行
将这些YAML文件提交并推送到GitHub后，可以在GitHub仓库的"Actions"选项卡中看到工作流的运行情况。测试和分析完成后，GitHub会展示测试结果和静态检查报告。