<template>
  <v-app>
    <v-toolbar app>
      <v-toolbar-title class="headline text-uppercase">
        <span>Leaderboard</span>
      </v-toolbar-title>
    </v-toolbar>

    <v-content>
      <v-container grid-list-md
                   text-xs-center>
        <v-layout row
                  wrap>
          <v-flex offset-md2
                  md8
                  offset-xs2
                  xs8>
            <template>
              <v-data-table :rows-per-page-items="[50, 100, 150]"
                            :headers="headers"
                            :items="items"
                            :pagination.sync="pagination"
                            class="elevation-1">
                <template slot="items"
                          slot-scope="props">
                  <td class="text-xs-center">
                    {{ (pagination.page - 1)* pagination.rowsPerPage + props.index + 1}}
                  </td>
                  <td class="text-xs-center">
                    <v-avatar slot="activator"
                              size="36px">
                      <img :src="`https://avatars.githubusercontent.com/${props.item.UserLogin}`"
                           alt="Avatar">
                    </v-avatar>
                  </td>
                  <td class="text-xs-center"><a :href="`https://github.com/${props.item.UserLogin}`"
                       target="_blank">{{ props.item.UserLogin }}</a></td>
                  <td class="text-xs-right">{{ props.item.IssueComments }}</td>
                  <td class="text-xs-right">{{ props.item.IssuesCreated }}</td>
                </template>
              </v-data-table>
            </template>

          </v-flex>

        </v-layout>
      </v-container>
    </v-content>
  </v-app>
</template>

<script>
import axios from "axios";

export default {
  name: "App",
  components: {},
  data() {
    return {
      pagination: {
        sortBy: 'IssueComments'
      },
      headers: [
        {
          text: "",
          value: "",
          align: "center",
        },
        {
          text: "",
          value: "",
          sortable: false
        },
        {
          text: "Login",
          value: "UserLogin",
          align: "center"
        },
        {
          text: "issue Comments",
          value: "IssueComments",
          align: "right"
        },
        {
          text: "issues opened",
          value: "IssuesCreated",
          align: "right"
        }
      ],
      items: []
      // items: [
      //   {
      //     UserID: 653013,
      //     UserLogin: "alexellisuk",
      //     IssueComments: 2,
      //     IssuesCreated: 1
      //   }
      // ]
    };
  },
  async mounted() {
    const result = await axios.get("/leaderboard");
    this.items = result.data;
  }
};
</script>
