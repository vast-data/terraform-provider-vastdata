#Create a user named example with uid of 9000
resource vastdata_user example-user {
  name = "example"
  uid = 9000
}

#Create a user named user1 with leading group & supplementary groups
resource vastdata_group group2 {
name = "group2"
gid = 2000
}

resource vastdata_group group4 {
name = "group4"
gid = 4000
}


resource vastdata_user user1 {
name = "user1"
uid = 3000
leading_gid = resource.vastdata_group.group1.gid
gids = [
  resource.vastdata_group.group2.gid,
  resource.vastdata_group.group4.gid
]

}
