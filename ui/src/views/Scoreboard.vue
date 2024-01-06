<script setup>
import { onMounted, ref } from "vue";
import { useUserStore } from "@/stores/user.js";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";

const listData = ref([]);
const user = useUserStore().user;
const router = useRouter();
const sendScore = () => {
  if (!user) {
    router.push({
      path: "/",
    });
    return ElMessage.error("请先输入昵称");
  }
  const message = {
    type: "UserClick",
    nickname: user,
  };
  window.$ws.sendMessage(message);
};
const goHome = () => {
    router.push("/");
  };

onMounted(() => {
  console.log(window.$ws);
  window.$ws.onMessage((message) => {
    listData.value = message;
  });
});
</script>

<template>
  <el-row :gutter="10" class="list-container">
    <el-col :span="4">
      <div class="button-area">
        <!-- Return Button -->
        <el-button @click="goHome">Return</el-button>
        <!-- Click Button -->
        <el-button @click="sendScore">Click!</el-button>
      </div>
    </el-col>
    <el-col :span="20">
      <div class="list-item">
        <template v-for="item in listData">
          <div class="item" v-if="item.ID">
            <span class="name"> {{ item.ID }}</span>
            <span class="score"> {{ item.Score }}</span>
          </div>
        </template>
      </div>
    </el-col>
  </el-row>
</template>
<style scoped>
.button-area {
  width: 100%;
  height: 100%;
  display: flex;
  justify-content: center;
  align-items: center;
}

.list-container {
  padding: 10px;
}

.list-item {
  width: 100%;
  overflow: hidden;
}

.item {
  width: 100%;
  border: 1px solid #dcdfe6;
  margin: 20px 0;
  background: #fff;
  padding: 15px 10px;
  border-radius: 5px;
  color: #606266;
  display: flex;
  justify-content: space-between;
  box-sizing: border-box;
}
</style>
