export interface User {
  id: string
  updated: string
  name: string
  lastSeen: string
}

export interface UserEntities {
  users: {
    [uuid: string]: User
  }
}
