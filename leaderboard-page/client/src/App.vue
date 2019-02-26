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
                            class="elevation-1">
                <template slot="items"
                          slot-scope="props">
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
      headers: [
        {
          text: "",
          value: ""
        },
        {
          text: "Login",
          value: "UserLogin"
        },
        {
          text: "issue Comments",
          value: "IssueComments"
        },
        {
          text: "issues opened",
          value: "IssuesCreated"
        },
        {
          text: "pull requests opened",
          value: "PullRequestsCreated"
        },
        {
          text: "review comments",
          value: "PRReviewComments"
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
