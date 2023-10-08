import Server from './server'
export default {
  getSSHHosts(): Promise<string[]> {
    return Server.get('/cshell/host/list')
  },
  getSSHHostsByTag(tag: string): Promise<string[]> {
    return Server.get(`/cshell/host/list/tag/${tag}`)
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
