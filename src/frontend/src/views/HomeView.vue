<script setup>
import {ref} from 'vue';
import {useRouter} from 'vue-router';
import {useUserStore} from '@/stores/user.js';
import { ElMessage } from 'element-plus'; // 引入ElMessage组件

const nickName = ref('');
const router = useRouter();
const setUser = useUserStore().setUser;
const joinService = () => {
  if (!nickName.value) { // 检查昵称是否为空
    return ElMessage.error('请先输入昵称'); // 如果为空，显示错误消息
  }
  const message = {
    'type': 'NewUser',
    'nickname': nickName.value
  };
  window.$ws.sendMessage(message);
  setUser(nickName.value);
  router.push({
    path:'/board'
  })
};
</script>

<template>
 <div class="home-view-container">
   <el-row type="flex" justify="center" align="middle" class="title-row">
     <el-col :span="24">
       <h1 class="title">Click大作战</h1>
     </el-col>
   </el-row>
   <el-row type="flex" justify="center" align="middle" class="input-row">
     <el-col :span="12">
       <el-input placeholder="请输入昵称" v-model="nickName" size="large">
         <template #append>
           <el-button type="primary" @click="joinService">加入</el-button>
         </template>
       </el-input>
     </el-col>
   </el-row>
 </div>
</template>

<style scoped>
.home-view-container{
  width:800px;
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
}

.title-row {
  margin-bottom: 20px;
}

.title {
  text-align: center;
  font-size: 2em;
  color: #409EFF;
}

.input-row {
  margin-top: 20px;
}
</style>