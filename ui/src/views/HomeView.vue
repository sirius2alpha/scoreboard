<script setup>
import {ref} from 'vue';
import {useRouter} from 'vue-router';
import {useUserStore} from '@/stores/user.js';
const nickName = ref('');
const router = useRouter();
const setUser = useUserStore().setUser;
const joinService = () => {
  const message = {
    'type': 'NewUser',
    'nickname': nickName.value
  };
  window.$ws.sendMessage(message);
  setUser(nickName.value);
  router.push({
    path:'/backend'
  })
};
</script>

<template>
 <div class="home-view-container">
   <el-input placeholder="请输入昵称" v-model="nickName" size="large">
     <template #append>
       <el-button @click="joinService">加入</el-button>
     </template>
   </el-input>
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
</style>
