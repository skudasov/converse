#### Converse
###### Telegram support bot
Install for development:
1. Install glide package manager
    ```
    curl https://glide.sh/get | sh
    or if OS X
    brew install glide
    ```
2. Install dependencies and docker-compose
    ```
    glide install
    pip install docker-compose
    ```
3. Create new bot token using Telegram Bot Father, insert into bot.yaml
4. Create support chat only for agents, insert invite link to *bot.yaml*
5. Fill your support Agents data in *bot.yaml*
   ```
    bot_token: ""
    bot_company_name: CORP
    delivery_ratelimit: 30 // telegram limits up to 30 messages/sec from one bot
    support_chatid: -297716605
    support_link: https://t.me/joinchat/... // Options->Manage group->Create invite link
    default_sla: 12 # default sla in hours
    agents:
        - name: "agent1"
          chatid: *your_userid_here*
        - name: "agent2"
          chatid: *your_userid_here*
    db:
        # postgres connection settings (container)
        host: "127.0.0.1"
        dbname: "tgsup"
        user: "tgsup"
        password: "tgsup"
        sslmode: "disable"
        # creates schema and migrate db, set true for initial setup, then false when restarting
        migrate: true
   ```

6. Run
   ```
   docker-compose build && docker-compose up
   ```

### Usage manual
#### User
On the first start conversation will be automatically created and you will be registered

![start](./content/start1.png)

Enter short description of your problem, then feel free describe it precisely and attach any files

When problem is solved, press [/current]() to end conversation
![conversation](./content/conversation.png)

After that you may open another conversation using [/active]() or select any resolved/closed conversation using [/history]()

![ending](./content/ending.png)

#### Agent
On the first start use invite link to join support chat

![start-agent](./content/start-agent.png)

From agent point of view you can participate in any conversation, one at a time, using [/active](), or review/reopen conversation from [/history]()

Use [/search]() **any text or caption of msg** when talking to bot to reveal helpful solutions from your history

![search](./content/search.png)

All conversation data stored in database before forwarding to Telegram, and will be restored when you join conversation, so in case of history glitch you will not lose any of your data

![restoration](./content/restore.png)

#### Support chat
In support chat one can see notifications of users and agents actions

![actions](./content/status.png)

Use [/list]() in *support chat* to see relevant conversations, ordered by time spent

![list](./content/list.png)
