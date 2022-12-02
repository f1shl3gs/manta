export interface User {
  id: string
  updated: string
  name: string
}

export interface UserEntities {
  users: {
    [uuid: string]: User
  }
}
