<template>
    <common-layout>
      <el-button type="primary" @click="save" style="right: 50px; top: 30px; position: absolute; z-index: 100">保存</el-button>
      <div style="padding: 15px 100px 10px 10px" :key="editHostId || 'new'">
        <el-button v-if="editHostId" type="danger" @click="confirmDel = true">删除配置</el-button>
        <el-form :model="config" label-width="150px">
          <el-form-item label="配置名">
            <el-input v-model="config.host"></el-input>
          </el-form-item>
          <el-form-item label="标签">
            <el-select v-model="config.tags" multiple allow-create filterable style="display: flex">
              <el-option v-for="tag in tags" :key="tag" :value="tag"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item label="可见">
            <el-switch v-model="config.visible"></el-switch>
          </el-form-item>
          <el-form-item label="主机名">
            <el-input v-model="config.f.HostName"></el-input>
          </el-form-item>
          <el-form-item label="端口">
            <el-input v-model="config.f.Port"></el-input>
          </el-form-item>
          <el-form-item label="用户名">
            <el-input v-model="config.f.User"></el-input>
          </el-form-item>
          <el-form-item label="使用钥匙串">
            <el-switch v-model="config.f.UseKeychain" active-value="yes" inactive-value="no"></el-switch>
          </el-form-item>
          <el-form-item label="AddKeysToAgent">
            <el-switch v-model="config.f.AddKeysToAgent" active-value="yes" inactive-value="no"></el-switch>
          </el-form-item>
          <el-form-item label="密钥文件">
            <el-select v-model="config.f.IdentityFile" filterable allow-create>
              <el-option v-for="k in keys" :value="k" :key="k"></el-option>
            </el-select>
          </el-form-item>
          <el-form-item>
            <template slot="label">
              <div style="display: flex; flex-direction: column">
                <span>代理主机</span>
                <span style="opacity: 0.5; font-size: 6px">ProxyJump</span>
              </div>
            </template>
            <el-input v-model="config.f.ProxyJump"></el-input>
          </el-form-item>
          <el-form-item>
            <template slot="label">
              <div style="display: flex; flex-direction: column">
                <span>代理命令</span>
                <span style="opacity: 0.5; font-size: 6px">ProxyCommand</span>
              </div>
            </template>
            <el-input v-model="config.f.ProxyCommand"></el-input>
          </el-form-item>
          <el-form-item label="额外配置">
            <el-row v-for="(_, idx) in addition" :key="idx">
              <el-col :span="21">
                <el-form-item label="Key" label-width="50px">
                  <el-input v-model="addition[idx].Key"></el-input>
                </el-form-item>
                <el-form-item label="Value" label-width="50px">
                  <el-input v-model="addition[idx].Value"></el-input>
                </el-form-item>
              </el-col>
              <el-col :offset="1" :span="2">
                <el-button type="danger" @click="addition.splice(idx, 1)">删除</el-button>
              </el-col>
            </el-row>
            <el-button type="primary" @click="addition.push({Key: '', Value: ''})">增加...</el-button>
          </el-form-item>
        </el-form>
      </div>
      <el-dialog :visible.sync="confirmDel" title="警告" modal>
        <span>确定要删除配置{{editHostId}}吗？</span>
        <span slot="footer" class="dialog-footer">
          <el-button @click="dialogVisible = false">取 消</el-button>
          <el-button type="danger" @click="deleteHost">确 定</el-button>
        </span>
      </el-dialog>
    </common-layout>
</template>

<script lang="ts">
import { Component, Vue, Watch } from 'vue-property-decorator'
import Service from '@/service'
import { getModule } from 'vuex-module-decorators'
import { Hosts } from '@/store/modules/hosts'

@Component({
    name: 'Editor',
    components: { }
})
export default class Editor extends Vue {
  config : any = { host: '', visible: false, tags: [], f: {
      User: '',
      HostName: '',
      UseKeychain: 'no',
      AddKeysToAgent: 'no',
      IdentityFile: '',
      Port: '',
      ProxyJump: '',
      ProxyCommand: ''
    } }
  addition: any[] = []
  get editHostId() {
    console.log(this.$route.params)
    return this.$route.params.id || ''
  }
  keys : string[] = []
  tags : string[] = []

  confirmDel = false

  specialKeys = ['User', 'HostName', 'UseKeychain', 'AddKeysToAgent', 'IdentityFile', 'Port', 'ProxyJump', 'ProxyCommand']
  special = new Set(this.specialKeys.map(v => v.toLowerCase()))

  @Watch('editHostId')
  async loadEditorInfo() {
    this.keys = await Service.getKeys()
    this.tags = await Service.getTags()
    this.config = { host: '', visible: false, tags: [], f: {
        User: '',
        HostName: '',
        UseKeychain: 'no',
        AddKeysToAgent: 'no',
        IdentityFile: '',
        Port: '',
        ProxyJump: '',
        ProxyCommand: ''
      } }
    if (!!this.editHostId) {
      const res : any = await Service.getHostConfig(this.editHostId)
      this.config.host = res.host as string
      this.config.visible = res.visible as boolean
      if (res.tags) {
        this.config.tags = res.tags
      }
      res.items.forEach((v: any) => {
        this.config.f[v.Key] = v.Value
      })
      this.addition = res.items.filter((v: any) => !this.special.has(v.Key.toLowerCase()))
    }
    console.log(this.config)
  }

  created() {
  }

  mounted() {
    this.loadEditorInfo()
  }

  genConfig() {
    const cfg = {
      host: this.config.host,
      visible: this.config.visible,
      tags: this.config.tags,
      items: [...this.addition],
    }
    this.specialKeys.forEach(v => {
      if (this.config.f[v]) {
        cfg.items.push({ Key: v, Value: this.config.f[v] })
      }
    })
    return cfg
  }

  async deleteHost() {
    await Service.deleteHostConfig(this.editHostId)
    await getModule(Hosts, this.$store).updateList()
    this.$nextTick(() => {
      this.$router.push('/')
    })
  }

  async save() {
    let hostId = this.editHostId
    if (hostId === '') {
      hostId = this.config.host
    }
    const cfg = this.genConfig()
    console.log(cfg)
    const res : any = await Service.setHostConfig(hostId, this.genConfig())
    this.$message({
      type: 'success',
      message: '保存成功。'
    })
    if (this.editHostId !== this.config.host) {
      await getModule(Hosts, this.$store).updateList()
      this.$nextTick(() => {
        this.$router.push(`/edit/${this.config.host}`)
      })
    }
  }
}
</script>

<style scoped lang="scss">
.el-row {
  margin-bottom: 20px;
  &:last-child {
    margin-bottom: 0;
  }
}
</style>
