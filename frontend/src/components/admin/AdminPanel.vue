<template>
  <v-container>
    <v-card>
      <v-card-title>Admin</v-card-title>
      <v-card-text>
        <v-alert v-if="!isSuperAdmin && !isTenantAdmin" type="warning" variant="tonal" class="mb-4">
          Admin access is required.
        </v-alert>
        <v-tabs v-model="tab" density="compact">
          <v-tab value="users">Users</v-tab>
          <v-tab value="tenants">Tenants</v-tab>
          <v-tab value="members">Members</v-tab>
          <v-tab value="accounts">AWS Accounts</v-tab>
          <v-tab value="audit">Audit</v-tab>
        </v-tabs>

        <v-window v-model="tab">
          <v-window-item value="users">
            <v-row class="mt-4" align="center">
              <v-col cols="12" md="6">
                <v-text-field v-model="userQuery" label="Search by email" variant="outlined" density="compact" />
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn color="primary" @click="loadUsers">Refresh</v-btn>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field v-model="newUser.email" label="New user email" variant="outlined" density="compact" />
              </v-col>
              <v-col cols="12" md="3">
                <v-text-field v-model="newUser.password" label="Password" type="password" variant="outlined" density="compact" />
              </v-col>
              <v-col cols="12" md="3">
                <v-checkbox v-model="newUser.isSuperAdmin" label="Super admin" density="compact" />
              </v-col>
              <v-col cols="12" md="3">
                <v-checkbox v-model="newUser.sendEmail" label="Send email" density="compact" />
              </v-col>
              <v-col cols="12" class="text-right">
                <v-btn color="primary" @click="createUser" :disabled="!isSuperAdmin">Create User</v-btn>
              </v-col>
            </v-row>
            <v-data-table :headers="userHeaders" :items="users" item-key="id" class="mt-4">
              <template v-slot:item.actions="{ item }">
                <v-btn size="small" variant="text" @click="toggleUserDisabled(item)" :disabled="!isSuperAdmin">
                  {{ item.isDisabled ? 'Enable' : 'Disable' }}
                </v-btn>
                <v-btn size="small" variant="text" @click="verifyUser(item)" :disabled="!isSuperAdmin || item.emailVerified">
                  Verify
                </v-btn>
                <v-btn size="small" variant="text" @click="resetPassword(item)" :disabled="!isSuperAdmin">
                  Reset Password
                </v-btn>
                <v-btn size="small" variant="text" color="error" @click="deleteUser(item)" :disabled="!isSuperAdmin">
                  Delete
                </v-btn>
              </template>
            </v-data-table>
          </v-window-item>

          <v-window-item value="tenants">
            <v-row class="mt-4" align="center">
              <v-col cols="12" md="8">
                <v-text-field v-model="newTenantName" label="New tenant name" variant="outlined" density="compact" />
              </v-col>
              <v-col cols="12" md="4" class="text-right">
                <v-btn color="primary" @click="createTenant" :disabled="!isSuperAdmin">Create Tenant</v-btn>
              </v-col>
            </v-row>
            <v-data-table :headers="tenantHeaders" :items="tenants" item-key="id" class="mt-4">
              <template v-slot:item.actions="{ item }">
                <v-btn size="small" variant="text" @click="renameTenant(item)" :disabled="!isSuperAdmin">Rename</v-btn>
                <v-btn size="small" variant="text" color="error" @click="deleteTenant(item)" :disabled="!isSuperAdmin">Delete</v-btn>
              </template>
            </v-data-table>
          </v-window-item>

          <v-window-item value="members">
            <v-row class="mt-4" align="center">
              <v-col cols="12" md="6" v-if="isSuperAdmin">
                <v-select
                  v-model="selectedTenantId"
                  :items="tenantSelectItems"
                  item-title="name"
                  item-value="id"
                  label="Select tenant"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn color="primary" @click="loadMembers">Refresh</v-btn>
              </v-col>
            </v-row>
            <v-row>
              <v-col cols="12" md="4">
                <v-text-field v-model="newMember.email" label="Member email" variant="outlined" density="compact" />
              </v-col>
              <v-col cols="12" md="2">
                <v-select
                  v-model="newMember.role"
                  :items="memberRoles"
                  label="Role"
                  variant="outlined"
                  density="compact"
                />
              </v-col>
              <v-col cols="12" md="3">
                <v-checkbox v-model="newMember.createIfMissing" label="Create if missing" density="compact" />
              </v-col>
              <v-col cols="12" md="3">
                <v-checkbox v-model="newMember.sendEmail" label="Send email" density="compact" />
              </v-col>
              <v-col cols="12" md="4">
                <v-text-field v-model="newMember.password" label="Temp password" type="password" variant="outlined" density="compact" />
              </v-col>
              <v-col cols="12" class="text-right">
                <v-btn color="primary" @click="addMember">Add Member</v-btn>
              </v-col>
            </v-row>
            <v-data-table :headers="memberHeaders" :items="members" item-key="userId" class="mt-4">
              <template v-slot:item.actions="{ item }">
                <v-btn size="small" variant="text" @click="changeMemberRole(item)">Change Role</v-btn>
                <v-btn size="small" variant="text" color="error" @click="removeMember(item)">Remove</v-btn>
              </template>
            </v-data-table>
          </v-window-item>

          <v-window-item value="accounts">
            <v-row class="mt-4" align="center">
              <v-col cols="12" class="text-right">
                <v-btn color="primary" @click="loadAwsAccounts">Refresh</v-btn>
              </v-col>
            </v-row>
            <v-data-table :headers="accountHeaders" :items="awsAccounts" item-key="id" class="mt-4" />
          </v-window-item>

          <v-window-item value="audit">
            <v-row class="mt-4" align="center">
              <v-col cols="12" md="6" v-if="isSuperAdmin">
                <v-select
                  v-model="auditTenantId"
                  :items="tenantSelectItems"
                  item-title="name"
                  item-value="id"
                  label="Filter tenant (optional)"
                  variant="outlined"
                  density="compact"
                  clearable
                />
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <v-btn color="primary" @click="loadAudit">Refresh</v-btn>
              </v-col>
            </v-row>
            <v-data-table :headers="auditHeaders" :items="auditLogs" item-key="id" class="mt-4" />
          </v-window-item>
        </v-window>
      </v-card-text>
    </v-card>
  </v-container>
</template>

<script>
import adminApi from '../../api/admin'

export default {
  name: 'AdminPanel',
  data() {
    return {
      tab: 'users',
      userQuery: '',
      users: [],
      tenants: [],
      members: [],
      awsAccounts: [],
      auditLogs: [],
      selectedTenantId: '',
      auditTenantId: '',
      newTenantName: '',
      newUser: {
        email: '',
        password: '',
        isSuperAdmin: false,
        sendEmail: false
      },
      newMember: {
        email: '',
        role: 'member',
        createIfMissing: false,
        sendEmail: false,
        password: ''
      },
      memberRoles: ['owner', 'admin', 'member'],
      userHeaders: [
        { title: 'Email', key: 'email' },
        { title: 'Verified', key: 'emailVerified' },
        { title: 'Super Admin', key: 'isSuperAdmin' },
        { title: 'Disabled', key: 'isDisabled' },
        { title: 'Created', key: 'createdAt' },
        { title: 'Actions', key: 'actions', sortable: false }
      ],
      tenantHeaders: [
        { title: 'Name', key: 'name' },
        { title: 'Members', key: 'members' },
        { title: 'Created', key: 'createdAt' },
        { title: 'Actions', key: 'actions', sortable: false }
      ],
      memberHeaders: [
        { title: 'Email', key: 'email' },
        { title: 'Role', key: 'role' },
        { title: 'Verified', key: 'emailVerified' },
        { title: 'Disabled', key: 'isDisabled' },
        { title: 'Actions', key: 'actions', sortable: false }
      ],
      accountHeaders: [
        { title: 'Tenant', key: 'tenantName' },
        { title: 'Account ID', key: 'accountId' },
        { title: 'Role ARN', key: 'roleArn' },
        { title: 'External ID', key: 'externalId' },
        { title: 'Name', key: 'name' },
        { title: 'Active', key: 'active' }
      ],
      auditHeaders: [
        { title: 'Actor', key: 'actorEmail' },
        { title: 'Action', key: 'action' },
        { title: 'Entity', key: 'entityType' },
        { title: 'Entity ID', key: 'entityId' },
        { title: 'Tenant', key: 'tenantId' },
        { title: 'Created', key: 'createdAt' }
      ]
    }
  },
  computed: {
    serverPath() {
      return this.$store.state.core.serverPath
    },
    isSuperAdmin() {
      return this.$store.state.auth.isSuperAdmin
    },
    isTenantAdmin() {
      return this.$store.state.auth.isTenantAdmin
    },
    currentTenantId() {
      return this.$store.state.auth.tenantId
    },
    tenantSelectItems() {
      return this.tenants || []
    }
  },
  mounted() {
    this.refreshAll()
  },
  methods: {
    refreshAll() {
      this.loadTenants()
      this.loadUsers()
      this.loadMembers()
      this.loadAwsAccounts()
      this.loadAudit()
    },
    loadUsers() {
      if (!this.isSuperAdmin) {
        this.users = []
        return
      }
      adminApi.listUsers(this.serverPath, { q: this.userQuery }, (resp) => {
        this.users = resp.data || []
      }, () => {})
    },
    createUser() {
      if (!this.isSuperAdmin) return
      adminApi.createUser(this.serverPath, this.newUser, () => {
        this.newUser = { email: '', password: '', isSuperAdmin: false, sendEmail: false }
        this.loadUsers()
      }, () => {})
    },
    toggleUserDisabled(user) {
      if (!this.isSuperAdmin) return
      const payload = {
        email: user.email,
        isSuperAdmin: user.isSuperAdmin,
        isDisabled: !user.isDisabled,
        emailVerified: user.emailVerified
      }
      adminApi.updateUser(this.serverPath, user.id, payload, () => {
        user.isDisabled = !user.isDisabled
      }, () => {})
    },
    verifyUser(user) {
      if (!this.isSuperAdmin) return
      adminApi.verifyEmail(this.serverPath, user.id, () => {
        user.emailVerified = true
      }, () => {})
    },
    resetPassword(user) {
      if (!this.isSuperAdmin) return
      const sendEmail = window.confirm('Send a reset email with a temporary password?')
      if (sendEmail) {
        adminApi.resetPassword(this.serverPath, user.id, { sendEmail: true }, () => {}, () => {})
        return
      }
      const password = window.prompt('Enter new password (min 8 chars)')
      if (!password) return
      adminApi.resetPassword(this.serverPath, user.id, { password }, () => {}, () => {})
    },
    deleteUser(user) {
      if (!this.isSuperAdmin) return
      if (!window.confirm(`Delete ${user.email}?`)) return
      adminApi.deleteUser(this.serverPath, user.id, () => {
        this.users = this.users.filter(u => u.id !== user.id)
      }, () => {})
    },
    loadTenants() {
      if (!this.isSuperAdmin) {
        this.tenants = []
        return
      }
      adminApi.listTenants(this.serverPath, (resp) => {
        this.tenants = resp.data || []
        if (!this.selectedTenantId && this.tenants.length) {
          this.selectedTenantId = this.tenants[0].id
        }
      }, () => {})
    },
    createTenant() {
      if (!this.isSuperAdmin) return
      if (!this.newTenantName) return
      adminApi.createTenant(this.serverPath, { name: this.newTenantName }, () => {
        this.newTenantName = ''
        this.loadTenants()
      }, () => {})
    },
    renameTenant(tenant) {
      if (!this.isSuperAdmin) return
      const name = window.prompt('New tenant name', tenant.name)
      if (!name) return
      adminApi.updateTenant(this.serverPath, tenant.id, { name }, () => {
        tenant.name = name
      }, () => {})
    },
    deleteTenant(tenant) {
      if (!this.isSuperAdmin) return
      if (!window.confirm(`Delete tenant ${tenant.name}?`)) return
      adminApi.deleteTenant(this.serverPath, tenant.id, () => {
        this.tenants = this.tenants.filter(t => t.id !== tenant.id)
      }, () => {})
    },
    loadMembers() {
      if (this.isSuperAdmin) {
        if (!this.selectedTenantId) {
          this.members = []
          return
        }
        adminApi.listTenantMembers(this.serverPath, this.selectedTenantId, (resp) => {
          this.members = resp.data || []
        }, () => {})
      } else if (this.isTenantAdmin) {
        adminApi.listTenantMembersSelf(this.serverPath, (resp) => {
          this.members = resp.data || []
        }, () => {})
      }
    },
    addMember() {
      const payload = { ...this.newMember }
      if (this.isSuperAdmin) {
        if (!this.selectedTenantId) return
        adminApi.addTenantMember(this.serverPath, this.selectedTenantId, payload, () => {
          this.newMember.email = ''
          this.newMember.password = ''
          this.loadMembers()
        }, () => {})
      } else if (this.isTenantAdmin) {
        adminApi.addTenantMemberSelf(this.serverPath, payload, () => {
          this.newMember.email = ''
          this.newMember.password = ''
          this.loadMembers()
        }, () => {})
      }
    },
    changeMemberRole(member) {
      const role = window.prompt('Role (owner/admin/member)', member.role)
      if (!role) return
      const payload = { role }
      if (this.isSuperAdmin) {
        adminApi.updateTenantMember(this.serverPath, this.selectedTenantId, member.userId, payload, () => {
          member.role = role
        }, () => {})
      } else if (this.isTenantAdmin) {
        adminApi.updateTenantMemberSelf(this.serverPath, member.userId, payload, () => {
          member.role = role
        }, () => {})
      }
    },
    removeMember(member) {
      if (!window.confirm(`Remove ${member.email} from tenant?`)) return
      if (this.isSuperAdmin) {
        adminApi.removeTenantMember(this.serverPath, this.selectedTenantId, member.userId, () => {
          this.members = this.members.filter(m => m.userId !== member.userId)
        }, () => {})
      } else if (this.isTenantAdmin) {
        adminApi.removeTenantMemberSelf(this.serverPath, member.userId, () => {
          this.members = this.members.filter(m => m.userId !== member.userId)
        }, () => {})
      }
    },
    loadAwsAccounts() {
      if (this.isSuperAdmin) {
        adminApi.listAwsAccounts(this.serverPath, (resp) => {
          this.awsAccounts = resp.data || []
        }, () => {})
      } else {
        adminApi.listAwsAccountsSelf(this.serverPath, (resp) => {
          this.awsAccounts = resp.data || []
        }, () => {})
      }
    },
    loadAudit() {
      const params = {}
      if (this.isSuperAdmin && this.auditTenantId) {
        params.tenant_id = this.auditTenantId
      }
      adminApi.listAuditLogs(this.serverPath, params, (resp) => {
        this.auditLogs = resp.data || []
      }, () => {})
    }
  }
}
</script>
