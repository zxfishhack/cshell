import Server from './server'
export default {
  getSSHHosts() {
    return Server.get('/cshell/host/list')
  },
  getHostConfig(hostId: string) {
    return Server.get(`/cshell/host/${hostId}/config`)
  },
  setHostConfig(hostId: string, config: any) {
    return Server.post(`/cshell/host/${hostId}/config`, config)
  },
  deleteHostConfig(hostId: string) {
    return Server.delete(`/cshell/host/${hostId}/config`)
  },
  getKeys() : Promise<string[]> {
    return Server.get('/cshell/keys')
  },
  getTags() : Promise<string[]> {
    return Server.get('/cshell/tags')
  }
}

export interface ServerInterface {
  [index: string]: Function
}
