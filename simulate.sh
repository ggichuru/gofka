#!/bin/bash

PROMPT="🤔"
DONE="🏆"
INFO="ℹ️"

menu_items=("🟢 Full Deployment", "🔵 Run Server", "EXIT🚪")

echo -e "🟢 Running simulation -> \t _ _ _ _ _ _ _\n"; date ;echo 

menu() {
    echo -e "\n${INFO} \tSimulate:"

    echo -ne "
    ${INFO} \tSimulation Menu
    1. Simulate [ User1 send to User2 ] 
    2. Simulate [ User1 receive from User2 ] 
    3. Simulate [ User1 receive from User4 ] 
    4. Get User Notifications 🔢
    5. Simulate All 🎮
    6. Start Kafka Broker  🐳
    0. Exit 🚪
    \nChoose an option:  ➕ " 

    read a
    case $a in
        1) u1sendU2 ; menu ;;
        2) u2sendU1 ; menu ;;
        3) u4sendU1 ; menu ;;
        4) retrieveNotifications ; menu ;;
        5) simulate_all ; menu ;;
        6) runKafka ; menu ;;
        0) exit 0 ;;
        *) echo -e "${RED}Wrong option.${STD}" && sleep 2
    esac

}


simulate_all() {
    echo -e "\nSimulate notificaitons\n"
    
}

u1sendU2() {
    # 
    echo -e "🎮 User One Sends User Two a Notification"
    curl -X POST http://localhost:8080/msgsend -d "fromID=1&toID=2&message=User2 mentioned you in a comment: 'Great seeing you yesterday, @User1!'"; echo


}

u2sendU1() {
    echo "🎮 User One Receives From User Two"
    curl -X POST http://localhost:8080/msgsend -d "fromID=2&toID=1&message=User1 started following you."; echo
}

u4sendU1() {
    echo "🎮 User One receives notifications from User Four"
    curl -X POST http://localhost:8080/msgsend -d "fromID=4&toID=1&message=Lena liked your post: 'My weekend getaway!'"; echo
}

retrieveNotifications() {
    echo -e "{PROMPT} Enter user ID"
    read input
    curl http://localhost:8081/notifications/$input
}

runKafka() {
    echo -e "🎮 Setting up Kafka Broker \n"
    docker-compose up -d
}


runNoticiationServers() {
    echo -e "🎮 Running notifications Producer \n"
    make test-producer && make test-consumer

    echo -e "🎮 Running notifications consumer \n"
}

# Call menu
menu