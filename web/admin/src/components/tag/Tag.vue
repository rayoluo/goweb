<template>
  <div>
    <a-card>
      <a-row :gutter="20">
        <a-col :span="4">
          <a-button type="primary" @click="addTagVisible = true">新增标签</a-button>
        </a-col>
      </a-row>

      <a-table
        rowKey="id"
        :columns="columns"
        :pagination="pagination"
        :dataSource="TagList"
        bordered
        @change="handleTableChange"
      >
        <template slot="name" slot-scope="name">
          <a-tag color="#87d068">{{ name }}</a-tag>
        </template>
        <template slot="action" slot-scope="data">
          <div class="actionSlot">
            <a-button type="primary" icon="edit" style="margin-right: 15px" @click="editTag(data.id)">编辑</a-button>
            <a-button type="danger" icon="delete" style="margin-right: 15px" @click="deleteTag(data.id)"
              >删除</a-button
            >
          </div>
        </template>
      </a-table>
    </a-card>

    <!-- 新增分类区域 -->
    <a-modal
      closable
      title="新增标签"
      :visible="addTagVisible"
      width="60%"
      @ok="addTagOk"
      @cancel="addTagCancel"
      destroyOnClose
    >
      <a-form-model :model="newTag" :rules="addTagRules" ref="addTagRef">
        <a-form-model-item label="标签名称" prop="name">
          <a-input v-model="newTag.name"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>

    <!-- 编辑分类区域 -->
    <a-modal
      closable
      destroyOnClose
      title="编辑标签"
      :visible="editTagVisible"
      width="60%"
      @ok="editTagOk"
      @cancel="editTagCancel"
    >
      <a-form-model :model="TagInfo" :rules="TagRules" ref="editTagRef">
        <a-form-model-item label="标签名称" prop="name">
          <a-input v-model="TagInfo.name"></a-input>
        </a-form-model-item>
      </a-form-model>
    </a-modal>
  </div>
</template>

<script>
const columns = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: '10%',
    key: 'id',
    align: 'center',
  },
  {
    title: '标签名',
    dataIndex: 'name',
    width: '20%',
    key: 'name',
    align: 'center',
    scopedSlots: { customRender: 'name' },
  },
  {
    title: '操作',
    width: '30%',
    key: 'action',
    align: 'center',
    scopedSlots: { customRender: 'action' },
  },
]

export default {
  data() {
    return {
      pagination: {
        pageSizeOptions: ['2', '5', '10', '20'],
        pageSize: 5,
        total: 0,
        showSizeChanger: true,
        showTotal: (total) => `共${total}条`,
      },
      TagList: [],
      TagInfo: {
        name: '',
        id: 0,
      },
      newTag: {
        name: '',
      },
      columns,
      queryParam: {
        pagesize: 5,
        pagenum: 1,
      },
      editVisible: false,
      TagRules: {
        name: [
          {
            validator: (rule, value, callback) => {
              if (this.TagInfo.name === '') {
                callback(new Error('请输入分类名'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      addTagRules: {
        name: [
          {
            validator: (rule, value, callback) => {
              if (this.newTag.name === '') {
                callback(new Error('请输入分类名'))
              } else {
                callback()
              }
            },
            trigger: 'blur',
          },
        ],
      },
      editTagVisible: false,
      addTagVisible: false,
    }
  },
  created() {
    this.getTagList()
  },
  methods: {
    // 获取标签列表
    async getTagList() {
      const { data: res } = await this.$http.get('admin/tag', {
        params: this.queryParam
      })

      if (res.status !== 200) {
        if (res.status === 1004 || res.status === 1005 || res.status === 1006 || res.status === 1007) {
          window.sessionStorage.clear()
          this.$router.push('/login')
        }
        this.$message.error(res.message)
      }
      this.TagList = res.data
      this.pagination.total = res.total
    },
    // 更改分页
    handleTableChange(pagination, filters, sorter) {
      var pager = { ...this.pagination }
      pager.current = pagination.current
      pager.pageSize = pagination.pageSize
      this.queryParam.pagesize = pagination.pageSize
      this.queryParam.pagenum = pagination.current

      if (pagination.pageSize !== this.pagination.pageSize) {
        this.queryParam.pagenum = 1
        pager.current = 1
      }
      this.pagination = pager
      this.getTagList()
    },
    // 删除标签
    deleteTag(id) {
      this.$confirm({
        title: '提示：请再次确认',
        content: '确定要删除该标签吗？一旦删除，无法恢复',
        onOk: async () => {
          const { data: res } = await this.$http.delete(`tag/${id}`)
          if (res.status != 200) return this.$message.error(res.message)
          // 删除标签时也应该删除标签和文章的对应关系
          const { data: res2 } = await this.$http.delete(`articleTag/tag/${id}`)
          if (res2.status !== 200) return this.$message.error('删除标签对应的标签文章对应关系失败!')
          this.$message.success('删除成功')
          this.getTagList()
        },
        onCancel: () => {
          this.$message.info('已取消删除')
        },
      })
    },
    // 新增分类
    addTagOk() {
      this.$refs.addTagRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const { data: res } = await this.$http.post('tag/add', {
          name: this.newTag.name,
        })
        if (res.status !== 200) return this.$message.error(res.message)
        this.$refs.addTagRef.resetFields()
        this.addTagVisible = false
        this.$message.success('添加标签成功')
        await this.getTagList()
      })
    },
    addTagCancel() {
      this.$refs.addTagRef.resetFields()
      this.addTagVisible = false
      this.$message.info('新增标签已取消')
    },
    // 编辑标签
    async editTag(id) {
      this.editTagVisible = true
      const { data: res } = await this.$http.get(`tag/${id}`)
      this.TagInfo = res.data
      this.TagInfo.id = id
    },
    editTagOk() {
      this.$refs.editTagRef.validate(async (valid) => {
        if (!valid) return this.$message.error('参数不符合要求，请重新输入')
        const { data: res } = await this.$http.put(`tag/${this.TagInfo.id}`, {
          name: this.TagInfo.name,
        })
        if (res.status != 200) return this.$message.error(res.message)
        this.editTagVisible = false
        this.$message.success('更新标签信息成功')
        this.getTagList()
      })
    },
    editTagCancel() {
      this.$refs.editTagRef.resetFields()
      this.editTagVisible = false
      this.$message.info('编辑已取消')
    },
  },
}
</script>

<style scoped>
.actionSlot {
  display: flex;
  justify-content: center;
}
</style>