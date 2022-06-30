<template>
  <div id="app">
    <el-tabs v-model="activeName" @tab-click="handleTabClick">
      <el-tab-pane label="实时消息" name="realtime">
        <div style="position: relative;z-index: 99;">
          <el-form :inline="true" @submit.native.prevent>
            <el-form-item>
              <el-input type="text" v-model="inputText"/>
            </el-form-item>
            <el-form-item>
              <el-button size="medium" type="primary" @click="sendBtnClick">发送</el-button>
            </el-form-item>
          </el-form>
        </div>
        <div>
          <el-table height="800" :data="msgList" stripe style="width: 100%;">
            <!--        <el-table-column prop="Sender" label="sender"></el-table-column>-->
            <!--        <el-table-column prop="Created" label="created"></el-table-column>-->
            <el-table-column prop="Body" label="message"></el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
      <el-tab-pane label="历史消息" name="history">
        <div>
          <el-table height="800" :data="historyMsgList" stripe style="width: 100%;">
            <!--        <el-table-column prop="Sender" label="sender"></el-table-column>-->
            <!--        <el-table-column prop="Created" label="created"></el-table-column>-->
            <el-table-column prop="Body" label="message"></el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="文件管理" name="fileList">
        <div>
          <el-upload
              ref="uploadbox"
              action="/api/v1/file/upload"
              :on-preview="handleFilePreview"
              :on-remove="handleFileRemove"
              :on-success="handleFileUploadSuccess"
              :file-list="uploadFileList"
              :auto-upload="false">
            <el-button slot="trigger" size="small" type="primary">选取文件</el-button>
            <el-button style="margin-left: 10px;" size="small" type="success" @click="submitFileUpload">上传到服务器</el-button>
          </el-upload>
        </div>
        <div>
          <el-table height="800" :data="fileList" stripe style="width: 100%;">
            <el-table-column prop="name" label="filename"></el-table-column>
            <el-table-column prop="size" label="size"></el-table-column>
            <el-table-column prop="created" label="created"></el-table-column>
            <el-table-column label="option">
              <template slot-scope="scope">
                <el-button type="danger" icon="el-icon-delete" size="small" @click="handleFileDelete(scope.row.name)"></el-button>
                <el-button type="success" icon="el-icon-download" size="small" @click="handleFileDownload(scope.row.name)"></el-button>
                <el-link :href="'/download/'+scope.row.name" :download="scope.row.name" :id="scope.row.name" v-show="false"></el-link>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="git管理" name="git-server">
        <el-popover width="400" trigger="click" v-model="gitServerPopoverVisible">
          <div>
            <el-form :inline="true" @submit.native.prevent>
              <el-form-item>
                <el-input type="text" v-model="newGitProjectName"></el-input>
              </el-form-item>
              <el-form-item>
                <el-button size="small" type="primary" @click="createNewGitProject">确定</el-button>
              </el-form-item>
              <el-form-item>
                <el-button size="small" @click="gitServerPopoverVisible = false">取消</el-button>
              </el-form-item>
            </el-form>
          </div>

          <el-button size="small" type="primary" slot="reference">创建新项目</el-button>
        </el-popover>
        <div>
          <el-table height="800" :data="gitProjectList" stripe style="width: 100%;">
            <el-table-column prop="name" label="name"></el-table-column>
            <el-table-column prop="path" label="uri"></el-table-column>
          </el-table>
        </div>
      </el-tab-pane>
    </el-tabs>

  </div>
</template>

<script>
const axios = require('axios');
const defaultChatUrl = `ws://${location.host}/api/v1/chatroom/online`;
export default {
  data: function() {
    return {
      serverInfo:{
        chatUrl: defaultChatUrl,
      },
      activeName: 'realtime',
      inputText: '',
      msgList: [],
      historyMsgList:[],
      fileList:[],
      uploadFileList:[],
      gitProjectList:[],
      newGitProjectName:'',
      gitServerPopoverVisible: false,
      ws: undefined,
    }
  },
  created() {
    this.serverInfo.chatUrl = "ws://" + location.host + "/api/v1/chatroom/online";
    this.connectServer();
  },
  methods: {
    connectServer() {
      this.ws = new WebSocket(this.serverInfo.chatUrl);
      this.ws.onopen = this.handleWsOnOpen;
      this.ws.onclose = this.handleWsOnClose;
      this.ws.onmessage = this.handleWsOnMessage;
      this.ws.onerror= this.handleWsOnError;
    },
    handleTabClick(tab, event) {
      if (tab.name === "history") {
        this.loadHistoryMsg(1, 20)
      }else if (tab.name === "fileList") {
        this.loadFileList()
      } else if (tab.name === "git-server"){
        this.freshGitProjectList();
      }else{
      }
    },
    loadHistoryMsg(pageNum, pageSize) {
      console.log("Load History");
      axios.post('/api/v1/chatroom/history',
          {
            "page_num": pageNum,
            "page_size": pageSize
          }
      ).then( response => {
        console.log(response);
        if (response.status === 200) {
          this.historyMsgList = response.data;
        }else{
          this.$message.error(response.statusText);
        }
      }).catch(error => {
        console.log(error);
      })
    },
    handleWsOnOpen(event) {
      this.$message.info("Connect to Server success");
    },
    handleWsOnClose(event) {
      this.ws = undefined;
    },
    handleWsOnMessage(event) {
      this.msgList.unshift(JSON.parse(event.data))
    },
    handleWsOnError(event) {
      this.$message.error("Connect to server failed");
      this.ws = undefined;
    },
    sendBtnClick() {
      if (undefined === this.ws || this.ws === null) {
        this.$message.error("Unconnect to server");
        return false;
      }
      this.ws.send(this.inputText);
    },
    freshGitProjectList() {
      axios.post("/api/v1/git/list", {})
      .then(resp => {
        this.gitProjectList = resp.data.body.list;
      })
      .catch(err => {
        this.$message.error(err)
      });
    },
    createNewGitProject() {
      if ('' !== this.newGitProjectName && undefined !== this.newGitProjectName && null !== this.newGitProjectName) {
        axios.post("/api/v1/git/add", {"name":this.newGitProjectName})
            .then(resp => {
              this.$message.info(resp.statusText);
            })
            .catch(err => {
              this.$message.error(err);
            });
      }
      this.gitServerPopoverVisible = false;
      this.freshGitProjectList();
    },
    loadFileList() {
      axios.post("/api/v1/file/list", {})
      .then(resp => {
        if (resp.status === 200) {
          console.log(resp);
          this.fileList = resp.data.body.files;
        }
      }).catch(err => {
        console.log(err);
      });
    },
    handleFileDelete(filename) {
     console.log(filename);
     axios.post("/api/v1/file/delete", {"filename":filename})
         .then(resp => {this.loadFileList();})
         .catch(err => {});
    },
    submitFileUpload() {
      this.$refs.uploadbox.submit();
    },
    handleFileRemove(file, fileList) {
      console.log(file, fileList);
    },
    handleFilePreview(file) {
      console.log(file);
    },
    handleFileUploadSuccess(resp) {
      if (resp.code === 0) {
        this.$message.info("upload success");
      }else{
        this.$message.error("upload failed" + resp.code_msg);
      }
      this.loadFileList();
    },
    handleFileDownload(filename) {
      // 直接调用a标签
      document.getElementById(filename).click();
    }
  }
}
</script>

<style>
#app {
  font-family: Helvetica, sans-serif;
  text-align: center;
}
</style>
