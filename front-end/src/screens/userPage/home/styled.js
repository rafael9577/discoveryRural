import {StyleSheet} from "react-native"

export default StyleSheet.create({
    safeArea: {
        flex: 1,
        justifyContent: 'center',
        alignItems: "center"
    },
    cardContainer: {
    width: '80%',
    height: '75%',
    borderWidth: 5,
    borderColor: '#2E7D32',
    borderRadius: 14,
    overflow: 'hidden', 
    alignSelf: 'center', 
    marginVertical: 'auto',
  },
  cardBackground: {
    width: '100%',
    height: '100%',
    justifyContent: 'flex-end', 
  },
  imageStyle: {
    borderRadius: 8, 
  },
  overlay: {
    backgroundColor: 'rgba(0, 0, 0, 0.4)', 
    padding: 16,
    width: '100%',
  },
  cardText: {
    color: '#FFFFFF',
    fontSize: 22,
    fontWeight: 'bold',
  },
    
})